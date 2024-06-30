package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	eventModels "github.com/BigNoseCattyHome/aorb/backend/services/event/models"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/services/poll/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
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
		}).Errorf("Error when marshal the event model")
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
		}).Errorf("Error when publishing the event model")
		logging.SetSpanError(span, err)
		return
	}
}

func (s PollServiceImpl) ListPolls(ctx context.Context, request *poll.ListPollRequest) (resp *poll.ListPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "ListVideosService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("FeedService.ListVideos").WithContext(ctx)

	now := time.Now().UnixMilli()
	latestTime := now
	if request.LatestTime != nil && *request.LatestTime != "" {
		// Check if request.LatestTime is a timestamp
		t, ok := isUnixMilliTimestamp(*request.LatestTime)
		if ok {
			latestTime = t
		} else {
			logger.WithFields(logrus.Fields{
				"latestTime": request.LatestTime,
			}).Errorf("The latestTime is not a unix timestamp")
			logging.SetSpanError(span, errors.New("the latestTime is not a unit timestamp"))
		}
	}

	find, nextTime, err := findPolls(ctx, latestTime)

	nextTimeStamp := uint64(nextTime.UnixMilli())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"find": find,
		}).Warnf("func findPolls meet trouble.")
		logging.SetSpanError(span, err)

		resp = &poll.ListPollResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			NextTime:   &nextTimeStamp,
			PollList:   nil,
		}
		return resp, err
	}

	if len(find) == 0 {
		resp = &poll.ListPollResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			NextTime:   nil,
			PollList:   nil,
		}
		return resp, nil
	}

	var actorId uint32 = 0
	if request.ActorId != nil {
		actorId = *request.ActorId
	}
	polls := queryDetailed(ctx, logger, actorId, find)
	if polls == nil {
		logger.WithFields(logrus.Fields{
			"polls": polls,
		}).Warnf("func queryDetailed meet trouble.")
		logging.SetSpanError(span, err)
		resp = &poll.ListPollResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			NextTime:   nil,
			PollList:   nil,
		}
		return resp, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var pollLists []uint32
		for _, item := range polls {
			pollLists = append(pollLists, item.Id)
		}
		producePoll(ctx, eventModels.RecommendEvent{
			ActorId: *request.ActorId,
			PollId:  pollLists,
			Type:    1,
			Source:  config.PollRpcServerName,
		})
	}()

	wg.Wait()
	resp = &poll.ListPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		NextTime:   &nextTimeStamp,
		PollList:   polls,
	}
	return resp, err
}

func (s PollServiceImpl) QueryPolls(ctx context.Context, req *poll.QueryPollRequest) (resp *poll.QueryPollResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "QueryPollsService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.QueryPolls").WithContext(ctx)

	rst, err := query(ctx, logger, req.ActorId, req.PollIds)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"rst": rst,
		}).Warnf("func query meet trouble.")
		logging.SetSpanError(span, err)
		resp = &poll.QueryPollResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			PollList:   rst,
		}
		return resp, err
	}

	resp = &poll.QueryPollResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		PollList:   rst,
	}
	return resp, err
}

func (s PollServiceImpl) QueryPollExisted(ctx context.Context, req *poll.PollExistRequest) (resp *poll.PollExistResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "QueryPollExistedService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("PollService.QueryPollExisted").WithContext(ctx)

	var tempPoll pollModels.Poll
	_, err = cached.GetWithFunc(ctx, fmt.Sprintf("PollExistedCached-%d", req.PollId), func(ctx context.Context, key string) (string, error) {
		collection := database.MongoDbClient.Database("aorb").Collection("polls")
		cursor := collection.FindOne(ctx, bson.M{"_id": req.PollId})
		if cursor.Err() != nil {
			return "false", cursor.Err()
		}
		if err := cursor.Decode(&tempPoll); err != nil {
			return "false", err
		}
		return "true", nil
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"poll_id": req.PollId,
		}).Warnf("Error occurred while querying database")
		logging.SetSpanError(span, err)
		resp = &poll.PollExistResponse{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			Existed:    false,
		}
		return resp, err
	}

	resp = &poll.PollExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Existed:    true,
	}
	return
}

func query(ctx context.Context, logger *logrus.Entry, actorId uint32, pollIds []uint32) (resp []*poll.Poll, err error) {
	var polls []*pollModels.Poll
	collection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.M{"_id": bson.M{"$in": pollIds}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warnf("Something went wrong when finding polls")
		return nil, err
	}
	for cursor.Next(ctx) {
		var tempPoll *pollModels.Poll
		if err = cursor.Decode(&tempPoll); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("Something went wrong when finding polls")
			return nil, err
		}
		polls = append(polls, tempPoll)
	}
	return queryDetailed(ctx, logger, actorId, polls), nil
}

func queryDetailed(ctx context.Context, logger *logrus.Entry, actorId uint32, polls []*pollModels.Poll) (respPollList []*poll.Poll) {
	ctx, span := tracing.Tracer.Start(ctx, "queryDetailed")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger = logging.LogService("ListVideos.queryDetailed").WithContext(ctx)
	respPollList = make([]*poll.Poll, len(polls))

	for i, v := range polls {
		respPollList[i] = &poll.Poll{
			Id:    v.ID,
			Title: v.Title,
			User: &user.User{
				Id: v.UserId,
			},
		}
	}

	userMap := make(map[uint32]*user.User)
	for _, poll := range polls {
		userMap[poll.UserId] = &user.User{}
	}

	userWg := sync.WaitGroup{}
	userWg.Add(len(userMap))
	for userId := range userMap {
		go func(userId uint32) {
			defer userWg.Done()
			userResponse, localErr := UserClient.GetUserInfo(ctx, &user.UserRequest{
				UserId:  userId,
				ActorId: actorId,
			})
			if localErr != nil || userResponse.StatusCode != strings.ServiceOKCode {
				logger.WithFields(logrus.Fields{
					"UserId": userId,
					"cause":  localErr,
				}).Warning("failed to get user info")
				logging.SetSpanError(span, localErr)
			}
			userMap[userId] = userResponse.User
		}(userId)
	}

	wg := sync.WaitGroup{}
	for i, p := range polls {
		wg.Add(1)
		go func(i int, v *pollModels.Poll) {
			defer wg.Done()
			commentCountResp, localErr := CommentClient.CountComment(ctx, &comment.CountCommentRequest{
				ActorId: actorId,
				PollId:  p.ID,
			})
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"poll_id": p.ID,
					"err":     localErr,
				}).Warning("failed to fetch comment count")
				logging.SetSpanError(span, localErr)
				return
			}
			respPollList[i].CommentCount = commentCountResp.CommentCount
		}(i, p)
	}

	userWg.Wait()
	wg.Wait()

	for i, respPoll := range respPollList {
		userId := respPoll.User.Id
		respPollList[i].User = userMap[userId]
	}
	return
}

func findPolls(ctx context.Context, latestTime int64) ([]*pollModels.Poll, time.Time, error) {
	logger := logging.LogService("ListPolls.findPolls").WithContext(ctx)

	nextTime := time.UnixMilli(latestTime)

	var polls []*pollModels.Poll
	collections := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.M{"created_at": bson.M{"$lte": nextTime}}
	sort := bson.M{"created_at": -1}
	findOptions := options.Find()
	findOptions.SetLimit(PollCount)
	findOptions.SetSort(sort)
	cursor, err := collections.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when finding the polls")
		return nil, nextTime, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var poll pollModels.Poll
		if err := cursor.Decode(&poll); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when finding the polls")
			return nil, nextTime, err
		}
		polls = append(polls, &poll)
	}

	if len(polls) != 0 {
		nextTime = polls[len(polls)-1].CreateAt
	}

	logger.WithFields(logrus.Fields{
		"nextTime":   nextTime,
		"latestTime": time.UnixMilli(latestTime),
		"PollsCount": len(polls),
	}).Debugf("find polls")
	return polls, nextTime, nil
}

func isUnixMilliTimestamp(s string) (int64, bool) {
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}

	startTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Now().AddDate(100, 0, 0)

	t := time.UnixMilli(timestamp)
	res := t.After(startTime) && t.Before(endTime)

	return timestamp, res
}
