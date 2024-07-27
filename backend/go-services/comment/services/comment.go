/*
service层, 负责处理业务逻辑
*/
package services

import (
	"context"
	"encoding/json"
	"fmt"
	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/event/models"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/BigNoseCattyHome/aorb/backend/utils/uuid"
	"github.com/go-redis/redis_rate/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"sync"
)

var userClient user.UserServiceClient
var pollClient poll.PollServiceClient
var actionCommentLimitKeyPrefix = config.Conf.Redis.Prefix + "comment_freq_limit"

const actionCommentMaxQPS = 3 // Maximum ActionComment query amount of an actor per second

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func actionCommentLimitKey(username string) string {
	return fmt.Sprintf("%s-%s", actionCommentLimitKeyPrefix, username)
}

type CommentServiceImpl struct {
	comment.CommentServiceServer
}

var conn *amqp.Connection
var channel *amqp.Channel

func (c CommentServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
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
		}).Errorf("Error with marshal the event models")
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
		}).Errorf("Error when publishing the event models")
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
		"username":     request.Username,
		"poll_id":      request.PollUuid,
		"action_type":  request.ActionType,
		"comment_text": request.GetCommentText(),
		"comment_uuid": request.GetCommentUuid(),
	}).Debugf("Process start")

	var pCommentId string
	var pCommentText string

	switch request.ActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		pCommentText = request.GetCommentText()
		break
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		pCommentId = request.GetCommentUuid()
		break
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
	limiterKey := actionCommentLimitKey(request.Username)
	limiterRes, err := limiter.Allow(ctx, limiterKey, redis_rate.PerSecond(actionCommentMaxQPS))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"userName": request.Username,
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
			"err":      err,
			"userName": request.Username,
		}).Infof("Action comment query too frequently by user %s", request.Username)

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.ActionCommentLimitedCode,
			StatusMsg:  strings.ActionCommentLimited,
		}
		return
	}

	// Check if poll exists
	pollExistResp, err := pollClient.PollExist(ctx, &poll.PollExistRequest{
		PollUuid: request.PollUuid,
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

	if !pollExistResp.Exist {
		logger.WithFields(logrus.Fields{
			"PollUuId": request.PollUuid,
		}).Errorf("Poll Uuid does not exist")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	// get target user
	userResponse, err := userClient.GetUserInfo(ctx, &user.UserRequest{
		Username: request.Username,
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
			"err":      err,
			"userName": request.Username,
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
		resp, err = addComment(ctx, logger, span, pUser, request.PollUuid, pCommentText)
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		resp, err = deleteComment(ctx, logger, span, pUser, request.PollUuid, pCommentId)
	}

	if err != nil {
		return
	}

	countCommentKey := fmt.Sprintf("Comment-Count-%s", request.PollUuid)
	cached.TagDelete(ctx, countCommentKey)

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debugf("Process done.")

	return
}

func (c CommentServiceImpl) ListComment(ctx context.Context, request *comment.ListCommentRequest) (resp *comment.ListCommentResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "ListCommentService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("CommentService.ListComment").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"poll_uuid": request.PollUuid,
	}).Debugf("Process start")

	// 查看poll是否存在
	pollExistResp, err := pollClient.PollExist(ctx, &poll.PollExistRequest{
		PollUuid: request.PollUuid,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query poll existence happens error")
		logging.SetSpanError(span, err)
		resp = &comment.ListCommentResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
		}
		return
	}

	if !pollExistResp.Exist {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
		}).Errorf("Poll ID does not exist")
		logging.SetSpanError(span, err)
		resp = &comment.ListCommentResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	var pPoll []pollModels.Poll
	collections := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{{"pollUuid", request.PollUuid}}
	cursor, err := collections.Find(ctx, filter)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": cursor.Err(),
		}).Errorf("CommentService list comment failed to response when listing comments")
		logging.SetSpanError(span, err)

		resp = &comment.ListCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}

	err = cursor.All(ctx, &pPoll)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error When Decoding Cursor of Poll")
		logging.SetSpanError(span, err)

		resp = &comment.ListCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}

	pCommentList := pPoll[0].CommentList

	// 获取每条评论的用户信息
	rCommentList := make([]*comment.Comment, 0, len(pCommentList))
	userMap := make(map[string]*user.User)
	for _, pComment := range pCommentList {
		userMap[pComment.ReviewerUserName] = &user.User{}
	}
	getUserInfoError := false

	wg := sync.WaitGroup{}
	wg.Add(len(userMap))
	for userName := range userMap {
		go func(userName string) {
			defer wg.Done()
			userResponse, getUserInfoErr := userClient.GetUserInfo(ctx, &user.UserRequest{
				Username: userName,
			})
			if err != nil || userResponse.StatusCode != strings.ServiceOKCode {
				logger.WithFields(logrus.Fields{
					"err":      getUserInfoErr,
					"userName": userName,
				}).Errorf("Unable to get user info")
				logging.SetSpanError(span, getUserInfoErr)
				getUserInfoError = true
				err = getUserInfoErr
			}
			userMap[userName] = userResponse.User
		}(userName)
	}
	wg.Wait()

	if getUserInfoError {
		resp = &comment.ListCommentResponse{
			StatusCode: strings.UnableToQueryUserErrorCode,
			StatusMsg:  strings.UnableToQueryUserError,
		}
		return
	}

	// 创建rCommentList
	for _, pComment := range pCommentList {
		curUser := userMap[pComment.ReviewerUserName]

		rCommentList = append(rCommentList, &comment.Comment{
			CommentUuid: pComment.CommentUuid,
			CommentUser: curUser,
			Content:     pComment.Content,
			CreateAt:    pComment.CreateAt,
		})
	}

	resp = &comment.ListCommentResponse{
		StatusCode:  strings.ServiceOKCode,
		StatusMsg:   strings.ServiceOK,
		CommentList: rCommentList,
	}

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debugf("Process done.")

	return
}

func (c CommentServiceImpl) CountComment(ctx context.Context, request *comment.CountCommentRequest) (resp *comment.CountCommentResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CountCommentService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("CommentService.CountComment").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"poll_uuid": request.PollUuid,
	}).Debugf("Process start")

	countStringKey := fmt.Sprintf("CommentCount-%s", request.PollUuid)
	countString, err := cached.GetWithFunc(ctx, countStringKey,
		func(ctx context.Context, key string) (string, error) {
			rCount, err := count(ctx, request.PollUuid)
			return strconv.FormatInt(rCount, 10), err
		})

	if err != nil {
		cached.TagDelete(ctx, "CommentCount")
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": request.PollUuid,
		}).Errorf("Unable to get comment count")
		logging.SetSpanError(span, err)

		resp = &comment.CountCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}
	rCount, _ := strconv.ParseUint(countString, 10, 64)
	resp = &comment.CountCommentResponse{
		StatusCode:   strings.ServiceOKCode,
		StatusMsg:    strings.ServiceOK,
		CommentCount: uint32(rCount),
	}
	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debugf("Process done.")
	return
}

func count(ctx context.Context, pollUuid string) (count int64, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CountComment")
	defer span.End()
	logger := logging.LogService("CommentService.CountComment").WithContext(ctx)

	collections := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{{"pollUuid", pollUuid}}
	result := collections.FindOne(ctx, filter)

	var pPoll pollModels.Poll
	result.Decode(&pPoll)
	count = int64(len(pPoll.CommentList))

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Faild to count comments")
		logging.SetSpanError(span, err)
	}
	return count, err
}

func deleteComment(ctx context.Context, logger *logrus.Entry, span trace.Span, pUser *user.User, pPollUuId string, commentUuid string) (resp *comment.ActionCommentResponse, err error) {
	rPoll := pollModels.Poll{}

	collections := database.MongoDbClient.Database("aorb").Collection("polls")

	filter := bson.D{
		{"pollUuid", pPollUuId},
		{"commentList.commentUuid", commentUuid},
	}

	// 先查询
	collections.FindOne(ctx, filter).Decode(&rPoll)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":          err,
			"poll_uuid":    pPollUuId,
			"comment_uuid": commentUuid,
		}).Errorf("Failed when searching for comment")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}

	flag := false

	for _, comment := range rPoll.CommentList {
		if comment.ReviewerUserName == pUser.Username {
			flag = true
			break
		}
	}

	if flag == false {
		// 没找到对应提问
		logger.WithFields(logrus.Fields{
			"err":          err,
			"poll_uuid":    pPollUuId,
			"comment_uuid": commentUuid,
		}).Errorf("Didn't find comment")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
		}
		return
	}

	filter = bson.D{{"pollUuid", pPollUuId}}
	update := bson.D{
		{"$pull", bson.D{
			{"commentList", bson.D{
				{"commentUuid", commentUuid},
			}},
		}},
	}
	_, err = collections.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":          err,
			"poll_uuid":    pPollUuId,
			"comment_uuid": commentUuid,
		}).Errorf("Failed to find comment")
		logging.SetSpanError(span, err)

		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToQueryCommentErrorCode,
			StatusMsg:  strings.UnableToQueryCommentError,
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

func addComment(ctx context.Context, logger *logrus.Entry, span trace.Span, pUser *user.User, pPollUuId string, pCommentText string) (resp *comment.ActionCommentResponse, err error) {

	collections := database.MongoDbClient.Database("aorb").Collection("polls")

	rComment := commentModels.Comment{
		CommentUuid:      uuid.GenerateUuid(),
		ReviewerUserName: pUser.Username,
		Content:          pCommentText,
		CreateAt:         timestamppb.Now(),
	}

	newComment := bson.D{
		{"commentUuid", rComment.CommentUuid},
		{"reviewerUserName", rComment.ReviewerUserName},
		{"content", rComment.Content},
		{"createAt", rComment.CreateAt},
	}

	update := bson.D{
		{"$push", bson.D{
			{"commentList", newComment},
		}},
	}

	filter := bson.D{
		{"pollUuid", pPollUuId},
	}

	_, err = collections.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": pPollUuId,
		}).Errorf("CommentService add comment action failed to response when adding comment")
		logging.SetSpanError(span, err)
		resp = &comment.ActionCommentResponse{
			StatusCode: strings.UnableToCreateCommentErrorCode,
			StatusMsg:  strings.UnableToCreateCommentError,
		}
		return
	}

	// TODO 一段意义不明的代码
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	//productComment(ctx, models.RecommendEvent{
	//	//	ActorId: pUser.Id,
	//	//	PollId:  []uint32{pPollId},
	//	//	Type:    2,
	//	//	Source:  config.CommentRpcServerName,
	//	//})
	//}()
	//wg.Wait()

	resp = &comment.ActionCommentResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Comment: &comment.Comment{
			CommentUuid: rComment.CommentUuid,
			CommentUser: pUser,
			Content:     rComment.Content,
			CreateAt:    rComment.CreateAt,
		},
	}
	return
}
