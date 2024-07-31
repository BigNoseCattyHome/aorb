package services

import (
	"context"
	"sync"
	"time"

	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	voteModels "github.com/BigNoseCattyHome/aorb/backend/go-services/vote/models"
	commentPb "github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	pollPb "github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	userPb "github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	votePb "github.com/BigNoseCattyHome/aorb/backend/rpc/vote"
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
)

type PollServiceImpl struct {
	pollPb.PollServiceServer
}

const (
	PollCount = 3
)

var UserClient userPb.UserServiceClient
var CommentClient commentPb.CommentServiceClient

var conn *amqp.Connection

var channel *amqp.Channel

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
func (s PollServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
	UserClient = userPb.NewUserServiceClient(userRpcConn)
	commentRpcConn := grpc2.Connect(config.CommentRpcServerName)
	CommentClient = commentPb.NewCommentServiceClient(commentRpcConn)

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

func (s PollServiceImpl) CreatePoll(ctx context.Context, request *pollPb.CreatePollRequest) (resp *pollPb.CreatePollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CreatePollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.CreatePoll").WithContext(ctx)

	var optionsCount = []uint32{0, 0}

	newPoll := &pollModels.Poll{
		PollUuid:     uuid.GenerateUuid(),
		PollType:     request.Poll.GetPollType(),
		Title:        request.Poll.GetTitle(),
		Options:      request.Poll.GetOptions(),
		OptionsCount: optionsCount,
		UserName:     request.Poll.GetUsername(),
		CommentList:  make([]commentModels.Comment, 0),
		VoteList:     make([]voteModels.Vote, 0),
		CreateAt:     time.Now(),
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
		resp = &pollPb.CreatePollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			PollUuid:   newPoll.PollUuid,
		}
		return resp, err
	}

	resp = &pollPb.CreatePollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollUuid:   newPoll.PollUuid,
	}
	return resp, err
}

func (s PollServiceImpl) GetPoll(ctx context.Context, request *pollPb.GetPollRequest) (resp *pollPb.GetPollResponse, err error) {
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
		resp = &pollPb.GetPollResponse{
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
		resp = &pollPb.GetPollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Poll:       nil,
		}
		return
	}

	resp = &pollPb.GetPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Poll:       BuildPollPbModel(&pPoll),
	}

	return
}

func (s PollServiceImpl) ListPoll(ctx context.Context, request *pollPb.ListPollRequest) (resp *pollPb.ListPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "ListPollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.ListPoll").WithContext(ctx)

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when listing polls")
		resp = &pollPb.ListPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
			PollList:   []*pollPb.Poll{},
		}
		return
	}

	var pAllPollList []pollModels.Poll
	err = cursor.All(ctx, &pAllPollList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when decoding polls")
		resp = &pollPb.ListPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
			PollList:   []*pollPb.Poll{},
		}
		return
	}

	var pPollList []pollModels.Poll
	for i := request.Limit * (request.Offset - 1); i < request.Limit*request.Offset && i < uint32(len(pAllPollList)); i++ {
		pPollList = append(pPollList, pAllPollList[i])
	}

	var rPollList []*pollPb.Poll
	for _, pPoll := range pPollList {
		rPollList = append(rPollList, BuildPollPbModel(&pPoll))
	}

	resp = &pollPb.ListPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollList:   rPollList,
	}
	return
}

func (s PollServiceImpl) PollExist(ctx context.Context, req *pollPb.PollExistRequest) (resp *pollPb.PollExistResponse, err error) {
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
		resp = &pollPb.PollExistResponse{
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
		resp = &pollPb.PollExistResponse{
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
		resp = &pollPb.PollExistResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Exist:      false,
		}
		return resp, err
	}

	resp = &pollPb.PollExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Exist:      true,
	}
	return
}

func BuildPollPbModel(poll *pollModels.Poll) *pollPb.Poll {

	rCommentList := make([]*commentPb.Comment, 0)
	rVoteList := make([]*votePb.Vote, 0)

	pCommentList := poll.CommentList
	pVoteList := poll.VoteList

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, pComment := range pCommentList {
			rComment := &commentPb.Comment{
				CommentUsername: pComment.CommentUserName,
				Content:         pComment.Content,
				CommentUuid:     pComment.CommentUuid,
				CreateAt:        timestamppb.New(pComment.CreateAt),
			}
			rCommentList = append(rCommentList, rComment)
		}
	}()

	go func() {
		defer wg.Done()
		for _, pVote := range pVoteList {
			rVote := &votePb.Vote{
				VoteUuid:     pVote.VoteUuid,
				VoteUsername: pVote.VoteUserName,
				Choice:       pVote.Choice,
				CreateAt:     timestamppb.New(pVote.CreateAt),
			}
			rVoteList = append(rVoteList, rVote)
		}
	}()

	wg.Wait()

	return &pollPb.Poll{
		PollUuid:     poll.PollUuid,
		Title:        poll.Title,
		Options:      poll.Options,
		OptionsCount: poll.OptionsCount,
		PollType:     poll.PollType,
		Username:     poll.UserName,
		CommentList:  rCommentList,
		VoteList:     rVoteList,
		CreateAt:     timestamppb.New(poll.CreateAt),
	}
}
