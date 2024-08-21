package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/messageQueue"
	"sync"
	"time"

	redisUtil "github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	// 先查有没有这个用户
	userResponse, err := UserClient.GetUserInfo(ctx, &userPb.UserRequest{
		Username: request.Poll.Username,
	})

	if err != nil || userResponse == nil || userResponse.StatusCode != strings.ServiceOKCode {
		if userResponse == nil || userResponse.StatusCode == strings.UserNotExistedCode {
			resp = &pollPb.CreatePollResponse{
				StatusCode: strings.UserNotExistedCode,
				StatusMsg:  strings.UserNotExisted,
			}
			return
		}
		logger.WithFields(logrus.Fields{
			"err":      err,
			"userName": request.Poll.Username,
		}).Errorf("Poll service error")
		logging.SetSpanError(span, err)
		resp = &pollPb.CreatePollResponse{
			StatusCode: strings.UnableToQueryUserErrorCode,
			StatusMsg:  strings.UnableToQueryUserError,
		}
		return
	}

	var optionsCount = []uint32{0, 0}

	newPoll := &pollModels.Poll{
		PollUuid:     uuid.GenerateUuid(),
		PollType:     request.Poll.GetPollType(),
		Title:        request.Poll.GetTitle(),
		Options:      request.Poll.GetOptions(),
		Content:      request.Poll.Content,
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

	// 将pollUuid加入user的PollAsk_List中
	userCollection := database.MongoDbClient.Database("aorb").Collection("users")
	filter4InsertPollUuid2PollAsk := bson.D{
		{"username", newPoll.UserName},
	}
	update4InsertPollUuid2PollAsk := bson.D{
		{"$push", bson.D{
			{"pollask.pollids", newPoll.PollUuid},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filter4InsertPollUuid2PollAsk, update4InsertPollUuid2PollAsk)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": newPoll.PollUuid,
			"username":  newPoll.UserName,
			"err":       err,
		}).Errorf("Error when inserting poll_uuid into user %s's pollask_list", request.Poll.Username)
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

func (s PollServiceImpl) GetPoll(ctx context.Context, request *pollPb.GetPollRequest) (response *pollPb.GetPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetPollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.GetPoll").WithContext(ctx)

	// 先去redis找
	redisKey := request.PollUuid
	// 从redis中获取数据
	redisResult, err := redisUtil.RedisCommentClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("获取poll失败, Error when getting data from redis")
		logging.SetSpanError(span, err)
		response = &pollPb.GetPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	if redisResult != "" {
		var rPoll *pollPb.Poll
		err = json.Unmarshal([]byte(redisResult), &rPoll)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("获取poll失败, Error when getting data from redis")
			logging.SetSpanError(span, err)
			response = &pollPb.GetPollResponse{
				StatusCode: strings.UnableToQueryPollErrorCode,
				StatusMsg:  strings.UnableToQueryPollError,
			}
			return
		}
		response = &pollPb.GetPollResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			Poll:       rPoll,
		}
		logger.WithFields(logrus.Fields{
			"response": response,
		}).Debugf("Process done.")
		return
	}

	// redis里面没有
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
		logging.SetSpanError(span, err)
		response = &pollPb.GetPollResponse{
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
		logging.SetSpanError(span, err)
		response = &pollPb.GetPollResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Poll:       nil,
		}
		return
	}

	jsonBytes, err := json.Marshal(BuildPollPbModel(&pPoll))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("获取poll失败, Error when marshalling pPoll")
		logging.SetSpanError(span, err)
		response = &pollPb.GetPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	// 将数据存入redis
	err = redisUtil.RedisPollClient.Set(ctx, redisKey, jsonBytes, time.Hour).Err()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("获取poll失败, Error when setting poll of uuid %s into redis", request.PollUuid)
		logging.SetSpanError(span, err)
		response = &pollPb.GetPollResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	response = &pollPb.GetPollResponse{
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
		logging.SetSpanError(span, err)
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
		logging.SetSpanError(span, err)
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

func (s PollServiceImpl) PollExist(ctx context.Context, request *pollPb.PollExistRequest) (response *pollPb.PollExistResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "PollExistedService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.PollExisted").WithContext(ctx)

	// 从redis里面查询
	redisKey := request.PollUuid
	redisResult, err := redisUtil.RedisCommentClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("调用pollExist失败, Error when getting data from redis")
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	if redisResult != "" {
		var rPoll *pollPb.Poll
		err = json.Unmarshal([]byte(redisResult), &rPoll)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("调用pollExist失败, Error when getting data from redis")
			logging.SetSpanError(span, err)
			response = &pollPb.PollExistResponse{
				StatusCode: strings.UnableToQueryPollErrorCode,
				StatusMsg:  strings.UnableToQueryPollError,
			}
			return
		}
		response = &pollPb.PollExistResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			Exist:      true,
		}
		logger.WithFields(logrus.Fields{
			"response": response,
		}).Debugf("Process done.")
		return
	}

	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	cursor := collection.FindOne(ctx, bson.M{"pollUuid": request.PollUuid})
	if cursor.Err() != nil {
		logger.WithFields(logrus.Fields{
			"err":       cursor.Err(),
			"poll_uuid": request.PollUuid,
		}).Errorf("调用pollExist失败, Error when checking if poll exists")
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
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
			"poll_uuid": request.PollUuid,
		}).Errorf("调用pollExist失败, Error when decoding poll")
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
			StatusCode: strings.PollServiceInnerErrorCode,
			StatusMsg:  strings.PollServiceInnerError,
			Exist:      false,
		}
		return response, err
	}

	if pPoll.PollUuid == "" {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
		}).Warnf("poll_uuid %s doesn't exist", request.PollUuid)
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			Exist:      false,
		}
		return response, err
	}

	jsonBytes, err := json.Marshal(BuildPollPbModel(&pPoll))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("调用pollExist失败, Error when marshalling pPoll")
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	// 将数据存入redis
	err = redisUtil.RedisPollClient.Set(ctx, redisKey, jsonBytes, time.Hour).Err()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("调用pollExist失败, Error when setting poll of uuid %s into redis", request.PollUuid)
		logging.SetSpanError(span, err)
		response = &pollPb.PollExistResponse{
			StatusCode: strings.UnableToQueryPollErrorCode,
			StatusMsg:  strings.UnableToQueryPollError,
		}
		return
	}

	response = &pollPb.PollExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Exist:      true,
	}
	return
}

func (s PollServiceImpl) FeedPoll(ctx context.Context, request *pollPb.FeedPollRequest) (response *pollPb.FeedPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "FeedPollService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.FeedPoll").WithContext(ctx)

	latestTime := time.Now()
	latestTimestamp := timestamppb.New(latestTime)

	// TODO 以后需要添加鉴权
	//////////

	redisKeys, err := redisUtil.RedisPollClient.Keys(ctx, "*").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("Error when getting keys from redisPollClient")
		logging.SetSpanError(span, err)
		response = &pollPb.FeedPollResponse{
			StatusCode: strings.PollServiceFeedErrorCode,
			StatusMsg:  strings.PollServiceFeedError,
		}
		return response, err
	}
	if len(redisKeys) > 10 {
		// 保留前10条提问
		redisKeys = redisKeys[:10]
	}

	var rPollList []*pollPb.Poll
	for _, redisKey := range redisKeys {
		redisResult, err := redisUtil.RedisPollClient.Get(ctx, redisKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			logger.WithFields(logrus.Fields{
				"username": request.Username,
			}).Errorf("获取提问流失败,Error when getting data from redisPollClient with key %s", redisKey)
			logging.SetSpanError(span, err)
			response = &pollPb.FeedPollResponse{
				StatusCode: strings.PollServiceFeedErrorCode,
				StatusMsg:  strings.PollServiceFeedError,
			}
			return response, err
		}
		if redisResult != "" {
			var rPoll pollPb.Poll
			err = json.Unmarshal([]byte(redisResult), &rPoll)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"username": request.Username,
				}).Errorf("获取提问流失败,Error when getting unmarshalling from redisPollClient with key %s", redisKey)
				logging.SetSpanError(span, err)
				response = &pollPb.FeedPollResponse{
					StatusCode: strings.PollServiceFeedErrorCode,
					StatusMsg:  strings.PollServiceFeedError,
				}
				return response, err
			}
			rPollList = append(rPollList, &rPoll)
		}
	}

	if len(rPollList) == 10 {
		// 一次拿10条数据，不够再从数据库拿
		response = &pollPb.FeedPollResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			PollList:   rPollList,
			NextTime:   latestTimestamp,
		}
	}

	// 从数据库获取剩余数据
	pollCollection := database.MongoDbClient.Database("aorb").Collection("polls")

	filter := bson.D{
		{"createAt", bson.D{{"$lt", latestTime}}},
		{"pollUuid", bson.D{{"$nin", redisKeys}}},
	}
	options := options.Find().SetSort(bson.D{{"createAt", 1}}).SetLimit(int64(10 - len(rPollList)))

	cursor, err := pollCollection.Find(ctx, filter, options)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("获取提问流失败, Error when getting poll from mongodb")
		logging.SetSpanError(span, err)
		response = &pollPb.FeedPollResponse{
			StatusCode: strings.PollServiceFeedErrorCode,
			StatusMsg:  strings.PollServiceFeedError,
		}
		return response, err
	}
	var pPollList []*pollModels.Poll
	err = cursor.All(ctx, &pPollList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("获取提问流失败, Error when decoding poll from mongodb")
		logging.SetSpanError(span, err)
		response = &pollPb.FeedPollResponse{
			StatusCode: strings.PollServiceFeedErrorCode,
			StatusMsg:  strings.PollServiceFeedError,
		}
		return response, err
	}
	for _, pPoll := range pPollList {
		rPollList = append(rPollList, BuildPollPbModel(pPoll))
	}

	var nextTime *timestamppb.Timestamp
	if len(pPollList) > 0 {
		nextTime = rPollList[len(rPollList)-1].CreateAt
	}

	// 将rPollList中的内容存入消息队列，异步写入redis
	for _, rPoll := range rPollList {
		body, _ := json.Marshal(&rPoll)
		err = rabbitmq.SendMessage2MQ(body, messageQueue.Poll2RedisQueue)
		if err != nil {
			//logger.WithFields(logrus.Fields{
			//	"username": request.Username,
			//}).Errorf("获取提问流失败, Error when sending message to rabbitmq")
			//logging.SetSpanError(span, err)
			//response = &pollPb.FeedPollResponse{
			//	StatusCode: strings.PollServiceFeedErrorCode,
			//	StatusMsg:  strings.PollServiceFeedError,
			//}
			//return response, err
		}
	}

	response = &pollPb.FeedPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollList:   rPollList,
		NextTime:   nextTime,
	}
	return
}

func (s PollServiceImpl) GetChoiceWithPollUuidAndUsername(ctx context.Context, request *pollPb.GetChoiceWithPollUuidAndUsernameRequest) (response *pollPb.GetChoiceWithPollUuidAndUsernameResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetChoiceWithPollUuidAndUsernameService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.GetChoiceWithPollUuidAndUsername").WithContext(ctx)

	// 先查有没有该user
	checkUserExistsResponse, err := UserClient.CheckUserExists(ctx, &userPb.UserExistRequest{
		Username: request.Username,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"username":  request.Username,
		}).Errorf("获取choice失败, Error when calling rpc CheckUserExists")
		logging.SetSpanError(span, err)
		response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
			StatusCode: strings.UnableToGetChoiceCode,
			StatusMsg:  strings.UnableToGetChoice,
		}
		return
	}

	if !checkUserExistsResponse.Existed {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"username":  request.Username,
		}).Errorf("获取choice失败, username doesn't exist")
		logging.SetSpanError(span, err)
		response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
			StatusCode: strings.UnableToGetChoiceCode,
			StatusMsg:  strings.UnableToGetChoice,
		}
		return
	}

	// 先redis查一下
	redisKey := request.PollUuid
	redisResult, err := redisUtil.RedisPollClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"username":  request.Username,
		}).Errorf("获取choice失败, Error when getting data from redisPollClient with key %s", redisKey)
		logging.SetSpanError(span, err)
		response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
			StatusCode: strings.UnableToGetChoiceCode,
			StatusMsg:  strings.UnableToGetChoice,
		}
		return
	}
	if redisResult != "" {
		// 有poll，查一下username
		var rPoll pollPb.Poll
		err = json.Unmarshal([]byte(redisResult), &rPoll)
		for _, rVote := range rPoll.VoteList {
			if rVote.VoteUsername == request.Username {
				response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
					StatusCode: strings.ServiceOKCode,
					StatusMsg:  strings.ServiceOK,
					VoteUuid:   rVote.VoteUuid,
					Choice:     rVote.Choice,
				}
			}
		}
	}

	// redis找不到就去数据库找
	getPollResponse, err := s.GetPoll(ctx, &pollPb.GetPollRequest{PollUuid: request.PollUuid})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_uuid": request.PollUuid,
			"username":  request.Username,
		}).Errorf("获取choice失败, Error when getting poll %s from mongodb", request.PollUuid)
		logging.SetSpanError(span, err)
		response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
			StatusCode: strings.UnableToGetChoiceCode,
			StatusMsg:  strings.UnableToGetChoice,
		}
		return
	}

	// 找到了就加入redis
	pollJson, _ := json.Marshal(&getPollResponse.Poll)
	redisUtil.RedisPollClient.Set(ctx, redisKey, pollJson, time.Hour)

	rVoteList := getPollResponse.Poll.VoteList
	for _, rVote := range rVoteList {
		if rVote.VoteUsername == request.Username {
			response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
				StatusCode: strings.ServiceOKCode,
				StatusMsg:  strings.ServiceOK,
				VoteUuid:   rVote.VoteUuid,
				Choice:     rVote.Choice,
			}
			return
		}
	}

	// 还没找到，返回false
	logger.WithFields(logrus.Fields{
		"poll_uuid": request.PollUuid,
		"username":  request.Username,
	}).Errorf("获取choice失败, Error when searching for %s's choice of poll_uuid: %s, cuz user didn't make a choice", request.Username, request.PollUuid)
	logging.SetSpanError(span, err)
	response = &pollPb.GetChoiceWithPollUuidAndUsernameResponse{
		StatusCode: strings.UnableToGetChoiceCode,
		StatusMsg:  strings.UnableToGetChoice,
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
		Content:      poll.Content,
		Options:      poll.Options,
		OptionsCount: poll.OptionsCount,
		PollType:     poll.PollType,
		Username:     poll.UserName,
		CommentList:  rCommentList,
		VoteList:     rVoteList,
		CreateAt:     timestamppb.New(poll.CreateAt),
	}
}

func PollMQ2Redis(ctx context.Context, rPoll *pollPb.Poll) error {
	pollJsonByte, err := json.Marshal(&rPoll)
	if err != nil {
		return err
	}
	err = redisUtil.RedisPollClient.Set(ctx, rPoll.PollUuid, pollJsonByte, time.Hour).Err()
	return err
}
