package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagePb "github.com/BigNoseCattyHome/aorb/backend/rpc/message"
	userPb "github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	redisUtil "github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/BigNoseCattyHome/aorb/backend/utils/uuid"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"

	commentModels "github.com/BigNoseCattyHome/aorb/backend/services/message/models"
	messageModels "github.com/BigNoseCattyHome/aorb/backend/services/message/models"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/services/message/models"
	voteModels "github.com/BigNoseCattyHome/aorb/backend/services/message/models"
)

var userClient userPb.UserServiceClient

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type MessageServiceImpl struct {
	messagePb.MessageServiceServer
}

var conn *amqp.Connection
var channel *amqp.Channel

func (m MessageServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
	userClient = userPb.NewUserServiceClient(userRpcConn)

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

func (m MessageServiceImpl) MessageChat(ctx context.Context, request *messagePb.MessageChatRequest) (response *messagePb.MessageChatResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "MessageService.MessageChat")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("MessageService.MessageChat").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"from_username": request.FromUsername,
		"to_username":   request.ToUsername,
	}).Debugf("Process start")

	fromUserExistResponse, err := userClient.CheckUserExists(ctx, &userPb.UserExistRequest{
		Username: request.FromUsername,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query user existence happens error")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	if !fromUserExistResponse.Existed {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("From user does not exist")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	toUserExistResponse, err := userClient.CheckUserExists(ctx, &userPb.UserExistRequest{
		Username: request.ToUsername,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query user existence happens error")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	if !toUserExistResponse.Existed {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("To user does not exist")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	// redis 操作
	redisKey := fmt.Sprintf("chat-messages:%s:%s", request.FromUsername, request.ToUsername)
	redisResult, err := redisUtil.RedisMessageClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when getting messageChat from redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	if redisResult != "" {
		// 有数据，直接返回
		rMessageChatList := make([]*messagePb.Message, 0)
		err = json.Unmarshal([]byte(redisResult), &rMessageChatList)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"from_username": request.FromUsername,
				"to_username":   request.ToUsername,
			}).Errorf("Error when unmarshalling messageChat from redis")
			logging.SetSpanError(span, err)
			response = &messagePb.MessageChatResponse{
				StatusCode: strings.UnableToQueryMessageErrorCode,
				StatusMsg:  strings.UnableToQueryMessageError,
			}
			return
		}
		response = &messagePb.MessageChatResponse{
			StatusCode:  strings.ServiceOKCode,
			StatusMsg:   strings.ServiceOK,
			MessageList: rMessageChatList,
		}
		return
	}

	// 查询并且返回结果
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")
	filter := bson.D{
		{Key: "fromUserName", Value: request.FromUsername},
		{Key: "toUserName", Value: request.ToUsername},
	}
	cursor, err := messageCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createAt", Value: 1}}))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when querying message")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	pMessageList := make([]*messageModels.Message, 0)
	err = cursor.All(ctx, &pMessageList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when decoding message")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	rMessageList := make([]*messagePb.Message, 0)
	for _, pMessage := range pMessageList {
		rMessageList = append(rMessageList, BuildMessagePbModel(pMessage))
	}

	// 存入redis
	jsonBytes, err := json.Marshal(&rMessageList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when marshalling messageChat to redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}
	err = redisUtil.RedisMessageClient.Set(ctx, redisKey, jsonBytes, time.Hour).Err()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when setting messageChat to redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageChatResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	response = &messagePb.MessageChatResponse{
		StatusCode:  strings.ServiceOKCode,
		StatusMsg:   strings.ServiceOK,
		MessageList: rMessageList,
	}
	return
}

func (m MessageServiceImpl) MessageAction(ctx context.Context, request *messagePb.MessageActionRequest) (response *messagePb.MessageActionResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "MessageService.MessageAction")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("MessageService.MessageAction").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"from_username": request.FromUsername,
		"to_username":   request.ToUsername,
	}).Debugf("Process start")

	fromUserExistResponse, err := userClient.CheckUserExists(ctx, &userPb.UserExistRequest{
		Username: request.FromUsername,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query user existence happens error")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	if !fromUserExistResponse.Existed {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("From user does not exist")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	toUserExistResponse, err := userClient.CheckUserExists(ctx, &userPb.UserExistRequest{
		Username: request.ToUsername,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Query user existence happens error")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	if !toUserExistResponse.Existed {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("To user does not exist")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	var pMessageUuid string
	var pMessageContent string

	switch request.ActionType {
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_ADD:
		pMessageContent = request.GetMessageContent()

	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_DELETE:
		pMessageUuid = request.GetMessageUuid()

	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_UNSPECIFIED:
		fallthrough
	default:
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("To user does not exist")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	var messageTypeString string
	switch request.MessageType {
	case messagePb.MessageType_MESSAGE_TYPE_FOLLOW:
		messageTypeString = "follow"
	case messagePb.MessageType_MESSAGE_TYPE_COMMENT:
		messageTypeString = "comment"
	case messagePb.MessageType_MESSAGE_TYPE_VOTE:
		messageTypeString = "vote"
	case messagePb.MessageType_MESSAGE_TYPE_CHAT:
		messageTypeString = "chat"
	}

	switch request.ActionType {
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_ADD:
		response, err = addMessage(ctx, logger, span, request.FromUsername, request.ToUsername, pMessageContent, messageTypeString)
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_DELETE:
		response, err = deleteMessage(ctx, logger, span, request.FromUsername, request.ToUsername, pMessageUuid)
	}

	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": request.FromUsername,
			"to_username":   request.ToUsername,
		}).Errorf("Error when adding or deleting message")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	logger.WithFields(logrus.Fields{
		"response": response,
	}).Debugf("Process done.")

	return
}

func (m MessageServiceImpl) GetUserMessage(ctx context.Context, request *messagePb.GetUserMessageRequest) (response *messagePb.GetUserMessageResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "MessageService.MessageAction")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("MessageService.MessageAction").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"username": request.Username,
	}).Debugf("Process start")

	rFollowMessageList := make([]*messagePb.FollowMessage, 0)
	rCommentReplyMessageList := make([]*messagePb.CommentReplyMessage, 0)
	rVoteMessageList := make([]*messagePb.VoteMessage, 0)
	messageUuidList := make([]string, 0)

	// 先去redis找
	redisKeys, err := redisUtil.RedisMessageClient.Keys(ctx, fmt.Sprintf("chat-message:*:%s", request.Username)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("获取用户接收消息失败, Error when getting keys from redisMessageClient")
		logging.SetSpanError(span, err)
		response = &messagePb.GetUserMessageResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}
	for _, redisKey := range redisKeys {
		redisResult, redisErr := redisUtil.RedisMessageClient.Get(ctx, redisKey).Result()
		if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
			logger.WithFields(logrus.Fields{
				"username": request.Username,
			}).Errorf("获取用户接收消息失败, Error when getting data from redisKey %s", redisKey)
			logging.SetSpanError(span, redisErr)
			response = &messagePb.GetUserMessageResponse{
				StatusCode: strings.UnableToQueryMessageErrorCode,
				StatusMsg:  strings.UnableToQueryMessageError,
			}
			return
		}
		if redisResult != "" {
			var rMessage messagePb.Message
			redisErr = json.Unmarshal([]byte(redisResult), &rMessage)
			if redisErr != nil {
				logger.WithFields(logrus.Fields{
					"username": request.Username,
				}).Errorf("获取用户接收消息失败, Error when unmarshalling data from redisKey %s", redisKey)
				logging.SetSpanError(span, redisErr)
				response = &messagePb.GetUserMessageResponse{
					StatusCode: strings.UnableToQueryMessageErrorCode,
					StatusMsg:  strings.UnableToQueryMessageError,
				}
				return
			}
			if rMessage.HasBeenRead == false {
				// 返回未读消息
				if rMessage.MessageType == messagePb.MessageType_MESSAGE_TYPE_FOLLOW {
					rFollowMessageList = append(rFollowMessageList, BuildFollowMessagePbWithMessagePb(&rMessage))
				} else if rMessage.MessageType == messagePb.MessageType_MESSAGE_TYPE_COMMENT {
					rCommentReplyMessageList = append(rCommentReplyMessageList, BuildCommentReplyMessagePbWithMessagePb(&rMessage))
				} else if rMessage.MessageType == messagePb.MessageType_MESSAGE_TYPE_VOTE {
					rVoteMessageList = append(rVoteMessageList, BuildVoteMessagePbWithMessagePb(&rMessage))
				}
				// 删掉这条key，防止缓存不一致，保证redis中的消息都是已读的
				redisUtil.RedisMessageClient.Del(ctx, redisKey)
				// 做个记录，防止重复插入
				messageUuidList = append(messageUuidList, rMessage.MessageUuid)
			}
		}
	}

	// 去mongodb找剩下的未读消息
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")
	filter := bson.D{
		{Key: "toUserName", Value: request.Username},
		{Key: "hasBeenRead", Value: false},
		{Key: "messageUuid", Value: bson.D{{Key: "$nin", Value: messageUuidList}}},
	}
	cursor, _ := messageCollection.Find(ctx, filter)
	pMessageList := make([]*messageModels.Message, 0)
	cursor.All(ctx, &pMessageList)
	for _, pMessage := range pMessageList {
		if pMessage.MessageType == "follow" {
			rFollowMessageList = append(rFollowMessageList, BuildFollowMessagePbWithMessageModel(pMessage))
		} else if pMessage.MessageType == "comment" {
			rCommentReplyMessageList = append(rCommentReplyMessageList, BuildCommentReplyMessagePbWithMessageModel(pMessage))
		} else if pMessage.MessageType == "vote" {
			rVoteMessageList = append(rVoteMessageList, BuildVoteMessagePbWithMessageModel(pMessage))
		}
	}

	response = &messagePb.GetUserMessageResponse{
		StatusCode:           strings.ServiceOKCode,
		StatusMsg:            strings.ServiceOK,
		FollowMessages:       rFollowMessageList,
		CommentReplyMessages: rCommentReplyMessageList,
		VoteMessages:         rVoteMessageList,
	}
	return
}

func (m MessageServiceImpl) MarkMessageStatus(ctx context.Context, request *messagePb.MarkMessageStatusRequest) (response *messagePb.MarkMessageStatusResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "MessageService.MessageAction")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("MessageService.MessageAction").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"message_uuid": request.MessageUuid,
	}).Debugf("Process start")

	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")

	switch request.Status {
	case messagePb.MessageStatus_MESSAGE_STATUS_READ:
		filter := bson.D{
			{Key: "messageUuid", Value: request.MessageUuid},
		}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "hasBeenRead", Value: true},
			}},
		}
		_, err = messageCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"message_uuid": request.MessageUuid,
			}).Errorf("修改阅读状态失败, Error when updating message status of message %s", request.MessageUuid)
			logging.SetSpanError(span, err)
			response = &messagePb.MarkMessageStatusResponse{
				StatusCode: strings.UnableToMarkMessageStatusCode,
				StatusMsg:  strings.UnableToMarkMessageStatus,
			}
			return
		}
		// 改完了加入redis
		var pMessage messageModels.Message
		cursor := messageCollection.FindOne(ctx, filter)
		cursor.Decode(&pMessage)
		rMessage := BuildMessagePbModel(&pMessage)
		redisKey := fmt.Sprintf("chat-message:%s:%s", rMessage.FromUsername, rMessage.ToUsername)
		jsonBytes, redisErr := json.Marshal(rMessage)
		if redisErr != nil {
			logger.WithFields(logrus.Fields{
				"message_uuid": request.MessageUuid,
			}).Errorf("修改阅读状态失败, Error when marshalling data from redisKey %s", redisKey)
			logging.SetSpanError(span, err)
			response = &messagePb.MarkMessageStatusResponse{
				StatusCode: strings.UnableToMarkMessageStatusCode,
				StatusMsg:  strings.UnableToMarkMessageStatus,
			}
			return
		}
		_, err = redisUtil.RedisMessageClient.Set(ctx, redisKey, jsonBytes, time.Hour).Result()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"message_uuid": request.MessageUuid,
			}).Errorf("修改阅读状态失败, Error when setting redisKey %s", redisKey)
			logging.SetSpanError(span, err)
			response = &messagePb.MarkMessageStatusResponse{
				StatusCode: strings.UnableToMarkMessageStatusCode,
				StatusMsg:  strings.UnableToMarkMessageStatus,
			}
			return
		}

	case messagePb.MessageStatus_MESSAGE_STATUS_UNREAD:
		filter := bson.D{
			{Key: "messageUuid", Value: request.MessageUuid},
		}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "hasBeenRead", Value: false},
			}},
		}
		_, err = messageCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"message_uuid": request.MessageUuid,
			}).Errorf("修改阅读状态失败, Error when updating message status of message %s", request.MessageUuid)
			logging.SetSpanError(span, err)
			response = &messagePb.MarkMessageStatusResponse{
				StatusCode: strings.UnableToMarkMessageStatusCode,
				StatusMsg:  strings.UnableToMarkMessageStatus,
			}
			return
		}
	}

	response = &messagePb.MarkMessageStatusResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}
	return
}

func BuildFollowMessagePbWithMessagePb(rMessage *messagePb.Message) *messagePb.FollowMessage {
	return &messagePb.FollowMessage{
		MessageUuid:      rMessage.MessageUuid,
		UsernameFollower: rMessage.FromUsername,
		Timestamp:        rMessage.CreateAt,
	}
}

func BuildFollowMessagePbWithMessageModel(rMessage *messageModels.Message) *messagePb.FollowMessage {
	return &messagePb.FollowMessage{
		MessageUuid:      rMessage.MessageUuid,
		UsernameFollower: rMessage.FromUserName,
		Timestamp:        timestamppb.New(rMessage.CreateAt),
	}
}

func BuildCommentReplyMessagePbWithMessagePb(rMessage *messagePb.Message) *messagePb.CommentReplyMessage {
	// 这里规定comment的message，其内容为comment_uuid
	// 先找到poll和comment
	pollCollection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{
		{Key: "commentList.commentUuid", Value: rMessage.Content},
	}
	cursor := pollCollection.FindOne(context.TODO(), filter)
	var pPoll pollModels.Poll
	cursor.Decode(&pPoll)

	var pCommentReply commentModels.Comment
	for _, pComment := range pPoll.CommentList {
		if pComment.CommentUserName == rMessage.FromUsername {
			pCommentReply = pComment
			break
		}
	}

	return &messagePb.CommentReplyMessage{
		MessageUuid: rMessage.MessageUuid,
		Content:     pCommentReply.Content,
		Username:    rMessage.FromUsername,
		PollUuid:    pPoll.PollUuid,
		CommentUuid: pCommentReply.CommentUuid,
		Timestamp:   rMessage.CreateAt,
	}
}

func BuildCommentReplyMessagePbWithMessageModel(rMessage *messageModels.Message) *messagePb.CommentReplyMessage {
	// 这里规定comment的message，其内容为comment_uuid
	// 先找到poll和comment
	pollCollection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{
		{Key: "commentList.commentUuid", Value: rMessage.Content},
	}
	cursor := pollCollection.FindOne(context.TODO(), filter)
	var pPoll pollModels.Poll
	cursor.Decode(&pPoll)

	var pCommentReply commentModels.Comment
	for _, pComment := range pPoll.CommentList {
		if pComment.CommentUserName == rMessage.FromUserName {
			pCommentReply = pComment
			break
		}
	}

	return &messagePb.CommentReplyMessage{
		MessageUuid: rMessage.MessageUuid,
		Content:     pCommentReply.Content,
		Username:    rMessage.FromUserName,
		PollUuid:    pPoll.PollUuid,
		CommentUuid: pCommentReply.CommentUuid,
		Timestamp:   timestamppb.New(rMessage.CreateAt),
	}
}

func BuildVoteMessagePbWithMessagePb(rMessage *messagePb.Message) *messagePb.VoteMessage {
	// 这里规定vote的message，其内容为vote_uuid
	// 先找到poll和vote
	pollCollection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{
		{Key: "commentList.voteUuid", Value: rMessage.Content},
	}
	cursor := pollCollection.FindOne(context.TODO(), filter)
	var pPoll pollModels.Poll
	cursor.Decode(&pPoll)

	var pVoteReply voteModels.Vote
	for _, pVote := range pPoll.VoteList {
		if pVote.VoteUserName == rMessage.FromUsername {
			pVoteReply = pVote
			break
		}
	}

	return &messagePb.VoteMessage{
		MessageUuid:  rMessage.MessageUuid,
		VoteUsername: rMessage.FromUsername,
		PollUuid:     pPoll.PollUuid,
		VoteUuid:     pVoteReply.VoteUuid,
		Timestamp:    rMessage.CreateAt,
		Choice:       pVoteReply.Choice,
	}
}

func BuildVoteMessagePbWithMessageModel(rMessage *messageModels.Message) *messagePb.VoteMessage {
	// 这里规定vote的message，其内容为vote_uuid
	// 先找到poll和vote
	pollCollection := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{
		{Key: "voteList.voteUuid", Value: rMessage.Content},
	}
	cursor := pollCollection.FindOne(context.TODO(), filter)
	var pPoll pollModels.Poll
	cursor.Decode(&pPoll)

	var pVoteReply voteModels.Vote
	for _, pVote := range pPoll.VoteList {
		if pVote.VoteUserName == rMessage.FromUserName {
			pVoteReply = pVote
			break
		}
	}

	return &messagePb.VoteMessage{
		MessageUuid:  rMessage.MessageUuid,
		VoteUsername: rMessage.FromUserName,
		PollUuid:     pPoll.PollUuid,
		VoteUuid:     pVoteReply.VoteUuid,
		Timestamp:    timestamppb.New(rMessage.CreateAt),
		Choice:       pVoteReply.Choice,
	}
}

func deleteMessage(ctx context.Context, logger *logrus.Entry, span trace.Span, fromUsername, toUserName string, messageUuid string) (response *messagePb.MessageActionResponse, err error) {
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")

	// 先查有没有message
	cursor := messageCollection.FindOne(ctx, bson.D{{Key: "messageUuid", Value: messageUuid}})

	if cursor == nil || cursor.Err() != nil {
		err = cursor.Err()
		logger.WithFields(logrus.Fields{
			"message_uuid": messageUuid,
		}).Errorf("Error when searching message_uuid: %s", messageUuid)
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	var pMessage messageModels.Message
	err = cursor.Decode(&pMessage)

	if err != nil {
		err = cursor.Err()
		logger.WithFields(logrus.Fields{
			"message_uuid": messageUuid,
		}).Errorf("Error when decoding message_uuid: %s", messageUuid)
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	_, err = messageCollection.DeleteOne(ctx, bson.D{{Key: "messageUuid", Value: messageUuid}})
	if err != nil {
		err = cursor.Err()
		logger.WithFields(logrus.Fields{
			"message_uuid": messageUuid,
		}).Errorf("Error when deleting message_uuid: %s", messageUuid)
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	// 删redis
	redisKey := fmt.Sprintf("chat-messages:%s:%s", fromUsername, toUserName)
	redisResult, err := redisUtil.RedisMessageClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		err = cursor.Err()
		logger.WithFields(logrus.Fields{
			"message_uuid": messageUuid,
		}).Errorf("Error when searching messageList from redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}
	if redisResult != "" {
		redisUtil.RedisMessageClient.Del(ctx, redisKey)
	}

	response = &messagePb.MessageActionResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}
	return
}

func addMessage(ctx context.Context, logger *logrus.Entry, span trace.Span, fromUserName, toUserName string, content string, messageType string) (response *messagePb.MessageActionResponse, err error) {
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")

	pMessage := &messageModels.Message{
		MessageUuid:  uuid.GenerateUuid(),
		FromUserName: fromUserName,
		ToUserName:   toUserName,
		Content:      content,
		HasBeenRead:  false,
		MessageType:  messageType,
		CreateAt:     time.Now(),
	}

	_, err = messageCollection.InsertOne(ctx, pMessage)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": fromUserName,
			"to_username":   toUserName,
		}).Errorf("Error when adding message")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	// 如果缓存中存在数据，则解码并合并到 messageList 中
	redisKey := fmt.Sprintf("chat-messages:%s:%s", fromUserName, toUserName)
	redisResult, err := redisUtil.RedisMessageClient.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.WithFields(logrus.Fields{
			"from_username": fromUserName,
			"to_username":   toUserName,
		}).Errorf("Error when getting messageList from redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	rActionMessageList := make([]*messagePb.Message, 0)
	if redisResult != "" {
		err = json.Unmarshal([]byte(redisResult), &rActionMessageList)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"from_username": fromUserName,
				"to_username":   toUserName,
			}).Errorf("Error when unmarshalling messageList from redis")
			logging.SetSpanError(span, err)
			response = &messagePb.MessageActionResponse{
				StatusCode: strings.UnableToQueryMessageErrorCode,
				StatusMsg:  strings.UnableToQueryMessageError,
			}
			return
		}
		rActionMessageList = append(rActionMessageList, BuildMessagePbModel(pMessage))
	}

	jsonBytes, err := json.Marshal(rActionMessageList)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": fromUserName,
			"to_username":   toUserName,
		}).Errorf("Error when marshalling messageList to redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}
	err = redisUtil.RedisMessageClient.Set(ctx, redisKey, jsonBytes, time.Hour).Err()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"from_username": fromUserName,
			"to_username":   toUserName,
		}).Errorf("Error when setting messageList to redis")
		logging.SetSpanError(span, err)
		response = &messagePb.MessageActionResponse{
			StatusCode: strings.UnableToQueryMessageErrorCode,
			StatusMsg:  strings.UnableToQueryMessageError,
		}
		return
	}

	response = &messagePb.MessageActionResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Message:    BuildMessagePbModel(pMessage),
	}
	return
}

func BuildMessagePbModel(pMessage *messageModels.Message) *messagePb.Message {
	var messageType int32
	switch pMessage.MessageType {
	case "follow":
		messageType = 0
		
	case "comment":
		messageType = 1
		
	case "vote":
		messageType = 2
		
	case "chat":
		messageType = 3
		
	}

	return &messagePb.Message{
		MessageUuid:  pMessage.MessageUuid,
		FromUsername: pMessage.FromUserName,
		ToUsername:   pMessage.ToUserName,
		Content:      pMessage.Content,
		MessageType:  messagePb.MessageType(messageType),
		HasBeenRead:  false,
		CreateAt:     timestamppb.New(pMessage.CreateAt),
	}
}
