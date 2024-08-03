package services

import (
	"context"
	"fmt"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	voteModels "github.com/BigNoseCattyHome/aorb/backend/go-services/vote/models"
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
	"time"
)

type VoteServiceImpl struct {
	votePb.VoteServiceServer
}

var userClient userPb.UserServiceClient
var pollClient pollPb.PollServiceClient

var conn *amqp.Connection

var channel *amqp.Channel

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func (s VoteServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
	userClient = userPb.NewUserServiceClient(userRpcConn)
	pollRpcConn := grpc2.Connect(config.PollRpcServerName)
	pollClient = pollPb.NewPollServiceClient(pollRpcConn)

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

func (s VoteServiceImpl) CreateVote(ctx context.Context, request *votePb.CreateVoteRequest) (response *votePb.CreateVoteResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "CreateVoteService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("VoteService.CreateVote").WithContext(ctx)

	// Check if poll exists
	pollExistResp, err := pollClient.PollExist(ctx, &pollPb.PollExistRequest{
		PollUuid: request.PollUuid,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query poll existence happens error")
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	if !pollExistResp.Exist {
		logger.WithFields(logrus.Fields{
			"PollUuId": request.PollUuid,
		}).Errorf("Poll Uuid does not exist")
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	// get target user
	userResponse, err := userClient.GetUserInfo(ctx, &userPb.UserRequest{
		Username: request.Username,
	})

	if err != nil || userResponse.StatusCode != strings.ServiceOKCode {
		if userResponse.StatusCode == strings.UserNotExistedCode {
			response = &votePb.CreateVoteResponse{
				StatusCode: strings.UserNotExistedCode,
				StatusMsg:  strings.UserNotExisted,
			}
			return
		}
		logger.WithFields(logrus.Fields{
			"err":      err,
			"userName": request.Username,
		}).Errorf("Vote service error")
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToQueryUserErrorCode,
			StatusMsg:  strings.UnableToQueryUserError,
		}
		return
	}

	collection := database.MongoDbClient.Database("aorb").Collection("polls")

	// Whether user had already voted or not
	filter4Check := bson.D{
		{"pollUuid", request.PollUuid},
	}
	var pPoll pollModels.Poll
	collection.FindOne(ctx, filter4Check).Decode(&pPoll)

	for _, vote := range pPoll.VoteList {
		if vote.VoteUserName == request.Username {
			logger.WithFields(logrus.Fields{
				"err":       err,
				"poll_uuid": request.PollUuid,
				"username":  request.Username,
			}).Errorf("user had already voted")
			logging.SetSpanError(span, err)
			response = &votePb.CreateVoteResponse{
				StatusCode: strings.UnableToCreateVoteErrorCode,
				StatusMsg:  strings.UnableToCreateVoteError,
			}
			return
		}
	}

	if request.Choice != pPoll.Options[0] && request.Choice != pPoll.Options[1] {
		logger.WithFields(logrus.Fields{
			"err":       "只能选择给定选项哦",
			"poll_uuid": request.PollUuid,
			"username":  request.Username,
		}).Errorf("options don't contain user's choice")
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToCreateVoteErrorCode,
			StatusMsg:  strings.UnableToCreateVoteError,
		}
		return
	}

	findIndex := func(options []string, choice string) int {
		for i, option := range options {
			if option == choice {
				return i
			}
		}
		return -1
	}

	pVote := &voteModels.Vote{
		VoteUuid:     uuid.GenerateUuid(),
		VoteUserName: request.Username,
		Choice:       request.Choice,
		CreateAt:     time.Now(),
	}

	filter := bson.D{{"pollUuid", request.PollUuid}}
	newVote := bson.D{
		{"voteUuid", pVote.VoteUuid},
		{"voteUserName", pVote.VoteUserName},
		{"choice", pVote.Choice},
		{"createAt", pVote.CreateAt},
	}
	update := bson.D{
		{"$push", bson.D{
			{"voteList", newVote},
		}},
		{"$inc", bson.D{
			{fmt.Sprintf("optionsCount.%d", findIndex(pPoll.Options, request.Choice)), 1},
		}},
	}

	_, err = collection.UpdateOne(ctx, filter, update)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": request.PollUuid,
		}).Errorf("VoteService create vote action failed to response when creating vote")
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	// 将对应的pollUuid加入user的pollans中
	userCollection := database.MongoDbClient.Database("aorb").Collection("users")
	filter4InsertPollUuid2PollAns := bson.D{
		{"username", pVote.VoteUserName},
	}
	update4InsertPollUuid2PollAns := bson.D{
		{"$push", bson.D{
			{"pollans.pollids", request.PollUuid},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filter4InsertPollUuid2PollAns, update4InsertPollUuid2PollAns)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": request.PollUuid,
			"username":  pVote.VoteUserName,
		}).Errorf("Error when inserting poll_uuid into user %s's pollans_list", pVote.VoteUserName)
		logging.SetSpanError(span, err)
		response = &votePb.CreateVoteResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	response = &votePb.CreateVoteResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		VoteUuid:   pVote.VoteUuid,
	}
	return
}

func (s VoteServiceImpl) GetVoteCount(ctx context.Context, request *votePb.GetVoteCountRequest) (response *votePb.GetVoteCountResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetVoteService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("VoteService.GetVote").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"poll_uuid": request.PollUuid,
	}).Debugf("Process start")

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{{"pollUuid", request.PollUuid}}
	cursor := collection.FindOne(ctx, filter)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": request.PollUuid,
		}).Errorf("Error when searching poll")
		logging.SetSpanError(span, err)
		response = &votePb.GetVoteCountResponse{
			StatusCode:    strings.UnableToQueryPollErrorCode,
			StatusMsg:     strings.UnableToQueryPollError,
			VoteCountList: []uint32{0, 0},
		}
		return
	}

	var pPoll pollModels.Poll
	err = cursor.Decode(&pPoll)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":       err,
			"poll_uuid": request.PollUuid,
		}).Errorf("Error when searching poll")
		logging.SetSpanError(span, err)
		response = &votePb.GetVoteCountResponse{
			StatusCode:    strings.UnableToGetVoteCountListErrorCode,
			StatusMsg:     strings.UnableToGetVoteCountListError,
			VoteCountList: []uint32{0, 0},
		}
		return
	}

	response = &votePb.GetVoteCountResponse{
		StatusCode:    strings.ServiceOKCode,
		StatusMsg:     strings.ServiceOK,
		VoteCountList: pPoll.OptionsCount,
	}
	return
}
