/*
service层, 负责处理业务逻辑
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/event/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/go-redis/redis_rate/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

var userClient user.UserServiceClient
var pollClient poll.PollServiceClient
var actionCommentLimitKeyPrefix = config.Conf.Redis.Prefix + "comment_freq_limit"
var rateCommentLimitKeyPrefix = config.Conf.Redis.Prefix + "comment_rate_limit"

const rateCommentMaxQPM = 3   // Maximum RateComment query amount
const actionCommentMaxQPS = 3 // Maximum ActionComment query amount of an actor per second

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func actionCommentLimitKey(userId string) string {
	return fmt.Sprintf("%s-%s", actionCommentLimitKeyPrefix, userId)
}

type CommentServiceImpl struct {
	comment.CommentServiceServer
}

var conn *amqp.Connection
var channel *amqp.Channel

func (c CommentServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServiceName)
	userClient = user.NewUserServiceClient(userRpcConn)

	var err error
	conn, err = amqp.Dial(rabbitmq.BuildMqConnAddr())
	exitOnError(err)

	channel, err = conn.Channel()
	exitOnError(err)

	err = channel.ExchangeDeclare(
		strings.EventExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	exitOnError(err)
}

func CloseMQConn() {
	if err := conn.Close(); err != nil {
		panic(err)
	}
	if err := channel.Close(); err != nil {
		panic(err)
	}
}

func productComment(ctx context.Context, event models.RecommendEvent) {
	ctx, span := tracing.Tracer.Start(ctx, "CommentPublisher")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("CommentService.CommentPublisher").WithContext(ctx)
	data, err := json.Marshal(event)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error with marshal the event model")
		logging.SetSpanError(span, err)
		return
	}

	headers := rabbitmq.InjectAMQPHeaders(ctx)

	err = channel.PublishWithContext(ctx,
		strings.EventExchange,
		strings.PollCommentEvent,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Headers:     headers,
		},
	)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when publishing the event model")
		logging.SetSpanError(span, err)
		return
	}
}

func (c CommentServiceImpl) ActionComment(ctx context.Context, request *comment.ActionCommentRequest) (resp *comment.ActionCommentResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CommentService.ActionComment")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("CommentService.ActionComment").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"user_id":      request.ActorId,
		"poll_id":      request.PollId,
		"action_type":  request.ActionType,
		"comment_text": request.GetCommentText(),
		"comment_id":   request.GetCommentId(),
	}).Debugf("Process start")

	var pCommentId string
	var pCommentText string

	switch request.ActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		pCommentText = request.GetCommentText()
		break
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		pCommentId = request.GetCommentId()
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		logger.Warnf("Invalid action type")
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.ActionCommentTypeInvalidCode,
			StatusMsg:  strings.ActionCommentTypeInvalid,
		}
		return
	}

	// Rate limiting
	limiter := redis_rate.NewLimiter(redis.Client)
	limiterKey := actionCommentLimitKey(request.ActorId)
	limiterRes, err := limiter.Allow(ctx, limiterKey, redis_rate.PerSecond(actionCommentMaxQPS))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"ActionId": request.ActorId,
		}).Errorf("ActionComment limiter error")
		logging.SetSpanError(span, err)

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToCreateCommentErrorCode,
			StatusMsg:  strings.UnableToCreateCommentError,
		}
		return
	}

	if limiterRes.Allowed == 0 {
		logger.WithFields(logrus.Fields{
			"err":     err,
			"ActorId": request.ActorId,
		}).Infof("Action comment query too frequently by user %d", request.ActorId)

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.ActionCommentLimitedCode,
			StatusMsg:  strings.ActionCommentLimited,
		}
		return
	}

	// Check if poll exists
	pollExistResp, err := pollClient.QueryPollExisted(ctx, &poll.PollExistRequest{
		PollId: request.PollId,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query video existence happens error")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
		}
		return
	}

	if !pollExistResp.Existed {
		logger.WithFields(logrus.Fields{
			"PollId": request.PollId,
		}).Errorf("Video ID does not exist")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryVideoErrorCode,
			StatusMsg:  strings.UnableToQueryVideoError,
		}
		return
	}

	// get target user
	userResponse, err := userClient.GetUserInfo(ctx, &user.UserRequest{
		UserId:  request.ActorId,
		ActorId: request.ActorId,
	})

	if err != nil || userResponse.StatusCode != strings.ServiceOKCode {
		if userResponse.StatusCode == strings.UserNotExistedCode {
			resp = &comment.ActionCommentResponse{
				StatusCode: strings.UserNotExistedCode,
				StatusMsg:  strings.UserNotExisted,
			}
			return
		}
		logger.WithFields(logrus.Fields{
			"err":     err,
			"ActorId": request.ActorId,
		}).Errorf("User service error")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryUserErrorCode,
			StatusMsg:  strings.UnableToQueryUserError,
		}
		return
	}

	pUser := userResponse.User

	switch request.ActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		resp, err = addComment(ctx, logger, span, pUser, request.PollId, pCommentText)
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		resp, err = deleteComment(ctx, logger, span, pUser, request.PollId, pCommentId)
	}

	if err != nil {
		return
	}

	countCommentKey := fmt.Sprintf("Comment-Count-%s", request.PollId)

}

func deleteComment(ctx context.Context, logger *logrus.Entry, span trace.Span, pUser *user.User, pPollId string, commentId string) (resp *comment.ActionCommentResponse, err error) {
	rComment := commentModels.Comment{}
	collections := database.MongoDbClient.Database("aorb").Collection("comments")
	result := collections.FindOne(ctx, rComment)
	if result.Err() != nil {
		logger.WithFields(logrus.Fields{
			"err":        result.Err(),
			"poll_id":    pPollId,
			"comment_id": commentId,
		}).Errorf("Failed to find comment")
		logging.SetSpanError(span, result.Err())

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}

	if rComment.UserId != pUser.Id {
		logger.Errorf("Comment creator and deletor not match")
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.ActorIDNotMatchErrorCode,
			StatusMsg:  strings.ActorIDNotMatchError,
		}
		return
	}

	_, err = collections.DeleteOne(ctx, rComment)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Failed to delete comment")
		logging.SetSpanError(span, err)

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToDeleteCommentErrorCode,
			StatusMsg:  strings.UnableToDeleteCommentError,
		}
		return
	}

	resp = &comment.ActionCommentResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Comment:    nil,
	}
	return
}

func addComment(ctx context.Context, logger *logrus.Entry, span trace.Span, pUser *user.User, pPollId string, pCommentText string) (resp *comment.ActionCommentResponse, err error) {
	rComment := commentModels.Comment{
		UserId:  pUser.Id,
		PollId:  pPollId,
		Content: pCommentText,
	}

	collections := database.MongoDbClient.Database("aorb").Collection("comments")
	_, err = collections.InsertOne(context.TODO(), rComment)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":        err,
			"comment_id": rComment.ID,
			"poll_id":    pPollId,
		}).Errorf("CommentService add comment action failed to response when adding comment")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToCreateCommentErrorCode,
			StatusMsg:  strings.UnableToCreateCommentError,
		}
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		productComment(ctx, models.RecommendEvent{
			ActorId: pUser.Id,
			PollId:  []string{pPollId},
			Type:    2,
			Source:  config.CommentRpcServiceName,
		})
	}()
	wg.Wait()

	resp = &comment.ActionCommentResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Comment: &comment.Comment{
			Id:       rComment.ID,
			User:     pUser,
			Content:  rComment.Content,
			CreateAt: rComment.CreateAt.Format("01-02"),
		},
	}
	return
}
