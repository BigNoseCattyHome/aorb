package services

import (
	"context"
	"encoding/json"
	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	eventModels "github.com/BigNoseCattyHome/aorb/backend/go-services/event/models"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	voteModels "github.com/BigNoseCattyHome/aorb/backend/go-services/vote/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/vote"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PollServiceImpl struct {
	poll.PollServiceServer
}

const (
	PollCount = 3
)

var UserClient user.UserServiceClient
var CommentClient comment.CommentServiceClient

var conn *amqp.Connection

var channel *amqp.Channel

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
func (s PollServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
	UserClient = user.NewUserServiceClient(userRpcConn)
	commentRpcConn := grpc2.Connect(config.CommentRpcServerName)
	CommentClient = comment.NewCommentServiceClient(commentRpcConn)

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

func producePoll(ctx context.Context, event eventModels.RecommendEvent) {
	ctx, span := tracing.Tracer.Start(ctx, "FeedPublisher")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("FeedService.FeedPublisher").WithContext(ctx)
	data, err := json.Marshal(event)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when marshal the event models")
		logging.SetSpanError(span, err)
		return
	}

	headers := rabbitmq.InjectAMQPHeaders(ctx)

	err = channel.PublishWithContext(ctx,
		strings.EventExchange,
		strings.PollGetEvent,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Headers:     headers,
		})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when publishing the event models")
		logging.SetSpanError(span, err)
		return
	}
}

func (s PollServiceImpl) CreatePoll(ctx context.Context, request *poll.CreatePollRequest) (resp *poll.CreatePollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CreatePollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.CreatePoll").WithContext(ctx)

	newPoll := pollModels.Poll{
		PollUuid:     uuid.GenerateUuid(),
		PollType:     request.Poll.PollType,
		Title:        request.Poll.Title,
		Options:      request.Poll.Options,
		OptionsCount: []int32{0, 0},
		UserName:     request.Poll.Username,
		CommentList:  []commentModels.Comment{},
		VoteList:     []voteModels.Vote{},
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	_, err = collection.InsertOne(ctx, newPoll)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": newPoll.PollUuid,
			"username":  newPoll.UserName,
			"err":       err,
		}).Errorf("Error when inserting new poll")
		logging.SetSpanError(span, err)
		resp = &poll.CreatePollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			PollUuid:   newPoll.PollUuid,
		}
		return resp, err
	}

	resp = &poll.CreatePollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollUuid:   newPoll.PollUuid,
	}
	return resp, err
}

func (s PollServiceImpl) GetPoll(ctx context.Context, request *poll.GetPollRequest) (resp *poll.GetPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetPollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.GetPoll").WithContext(ctx)

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{
		{"pollUuid", request.PollUuid},
	}

	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"err":       result.Err(),
		}).Errorf("Error when getting poll of uuid %s", request.PollUuid)
		resp = &poll.GetPollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Poll:       nil,
		}
		err = result.Err()
		return
	}

	var pPoll pollModels.Poll
	err = result.Decode(&pPoll)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"err":       err,
		}).Errorf("Error when decoding poll of uuid %s", request.PollUuid)
		resp = &poll.GetPollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Poll:       nil,
		}
		return
	}

	pCommentList := pPoll.CommentList
	var rCommentList []*comment.Comment
	for _, pComment := range pCommentList {
		rComment := comment.Comment{
			CommentUsername: pComment.ReviewerUserName,
			Content:         pComment.Content,
			CommentUuid:     pComment.CommentUuid,
			CreateAt:        timestamppb.New(pComment.CreateAt),
			DeleteAt:        nil,
		}
		rCommentList = append(rCommentList, &rComment)
	}

	pVoteList := pPoll.VoteList
	var rVoteList []*vote.Vote
	for _, pVote := range pVoteList {
		rVote := vote.Vote{
			VoteUuid:     pVote.VoteUuid,
			VoteUsername: pVote.VoteUserName,
			CreateAt:     timestamppb.New(pVote.CreateAt),
			DeleteAt:     timestamppb.New(pVote.DeleteAt),
		}
		rVoteList = append(rVoteList, &rVote)
	}

	resp = &poll.GetPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Poll: &poll.Poll{
			PollUuid:     pPoll.PollUuid,
			Title:        pPoll.Title,
			Options:      pPoll.Options,
			OptionsCount: pPoll.OptionsCount,
			PollType:     pPoll.PollType,
			Username:     pPoll.UserName,
			CommentList:  rCommentList,
			VoteList:     rVoteList,
			CreateAt:     timestamppb.New(pPoll.CreateAt),
			UpdateAt:     timestamppb.New(pPoll.UpdateAt),
			DeleteAt:     timestamppb.New(pPoll.DeleteAt),
		},
	}
	return
}

func (s PollServiceImpl) ListPoll(ctx context.Context, request *poll.ListPollRequest) (resp *poll.ListPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "ListVideosService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("FeedService.ListVideos").WithContext(ctx)

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when listing polls")
		resp = &poll.ListPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
			PollList:   []*poll.Poll{},
		}
		return
	}

	var pAllPollList []pollModels.Poll
	err = cursor.All(ctx, &pAllPollList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when decoding polls")
		resp = &poll.ListPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
			PollList:   []*poll.Poll{},
		}
		return
	}

	var pPollList []pollModels.Poll
	for i := request.Limit * (request.Offset - 1); i < request.Limit*request.Offset && i < int32(len(pAllPollList)); i++ {
		pPollList = append(pPollList, pAllPollList[i])
	}

	var rPollList []*poll.Poll
	for _, pPoll := range pPollList {
		var rCommentList []*comment.Comment
		for _, pComment := range pPoll.CommentList {
			rComment := &comment.Comment{
				CommentUsername: pComment.ReviewerUserName,
				Content:         pComment.Content,
				CommentUuid:     pComment.CommentUuid,
				CreateAt:        timestamppb.New(pComment.CreateAt),
				DeleteAt:        timestamppb.New(pComment.DeleteAt),
			}
			rCommentList = append(rCommentList, rComment)
		}

		var rVoteList []*vote.Vote
		for _, pVote := range pPoll.VoteList {
			rVote := &vote.Vote{
				VoteUuid:     pVote.VoteUuid,
				VoteUsername: pVote.VoteUserName,
				Choice:       pVote.Choice,
				CreateAt:     timestamppb.New(pVote.CreateAt),
				DeleteAt:     timestamppb.New(pVote.DeleteAt),
			}
			rVoteList = append(rVoteList, rVote)
		}

		rPoll := &poll.Poll{
			PollUuid:     pPoll.PollUuid,
			Title:        pPoll.Title,
			Options:      pPoll.Options,
			OptionsCount: pPoll.OptionsCount,
			PollType:     pPoll.PollType,
			Username:     pPoll.UserName,
			CommentList:  rCommentList,
			VoteList:     rVoteList,
			CreateAt:     timestamppb.New(pPoll.CreateAt),
			UpdateAt:     timestamppb.New(pPoll.UpdateAt),
			DeleteAt:     timestamppb.New(pPoll.DeleteAt),
		}

		rPollList = append(rPollList, rPoll)
	}

	resp = &poll.ListPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollList:   rPollList,
	}
	return
}

func (s PollServiceImpl) PollExist(ctx context.Context, req *poll.PollExistRequest) (resp *poll.PollExistResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "PollExistedService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.PollExisted").WithContext(ctx)

	// TODO 使用二级缓存加速查询
	//var tempPoll pollModels.Poll
	//_, err = cached.GetWithFunc(ctx, fmt.Sprintf("PollExistedCached-%s", req.PollUuid), func(ctx context.Context, key string) (string, error) {
	//	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	//	cursor := collection.FindOne(ctx, bson.M{"pollUuid": req.PollUuid})
	//	if cursor.Err() != nil {
	//		return "false", cursor.Err()
	//	}
	//	if err := cursor.Decode(&tempPoll); err != nil {
	//		return "false", err
	//	}
	//	return "true", nil
	//})

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	cursor := collection.FindOne(ctx, bson.M{"pollUuid": req.PollUuid})
	if cursor.Err() != nil {
		logger.WithFields(logrus.Fields{
			"err":       cursor.Err(),
			"poll_uuid": req.PollUuid,
		}).Errorf("Error when checking if poll exists")
		resp = &poll.PollExistResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
			Exist:      false,
		}
		return
	}

	var pPoll pollModels.Poll
	err = cursor.Decode(&pPoll)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": req.PollUuid,
		}).Errorf("Error when decoding poll")
		logging.SetSpanError(span, err)
		resp = &poll.PollExistResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Exist:      false,
		}
		return resp, err
	}

	if pPoll.PollUuid == "" {
		logger.WithFields(logrus.Fields{
			"poll_uuid": req.PollUuid,
		}).Warnf("poll_uuid %s doesn't exist", req.PollUuid)
		logging.SetSpanError(span, err)
		resp = &poll.PollExistResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Exist:      false,
		}
		return resp, err
	}

	resp = &poll.PollExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Exist:      true,
	}
	return
}
