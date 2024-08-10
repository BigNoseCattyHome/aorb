package services

import (
	"context"
	"errors"
	"fmt"
	messageModels "github.com/BigNoseCattyHome/aorb/backend/go-services/message/models"
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
	"time"
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
		{"fromUserName", request.FromUsername},
		{"toUserName", request.ToUsername},
	}
	cursor, err := messageCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{"createAt", 1}}))
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
		break
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_DELETE:
		pMessageUuid = request.GetMessageUuid()
		break
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

	switch request.ActionType {
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_ADD:
		response, err = addMessage(ctx, logger, span, request.FromUsername, request.ToUsername, pMessageContent)
		break
	case messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_DELETE:
		response, err = deleteMessage(ctx, logger, span, request.FromUsername, request.ToUsername, pMessageUuid)
		break
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

func deleteMessage(ctx context.Context, logger *logrus.Entry, span trace.Span, fromUsername, toUserName string, messageUuid string) (response *messagePb.MessageActionResponse, err error) {
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")

	// 先查有没有message
	cursor := messageCollection.FindOne(ctx, bson.D{{"messageUuid", messageUuid}})

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

	_, err = messageCollection.DeleteOne(ctx, bson.D{{"messageUuid", messageUuid}})
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

func addMessage(ctx context.Context, logger *logrus.Entry, span trace.Span, fromUserName, toUserName string, content string) (response *messagePb.MessageActionResponse, err error) {
	messageCollection := database.MongoDbClient.Database("aorb").Collection("messages")

	pMessage := &messageModels.Message{
		MessageUuid:  uuid.GenerateUuid(),
		FromUserName: fromUserName,
		ToUserName:   toUserName,
		Content:      content,
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
	return &messagePb.Message{
		MessageUuid:  pMessage.MessageUuid,
		FromUsername: pMessage.FromUserName,
		ToUsername:   pMessage.ToUserName,
		Content:      pMessage.Content,
		CreateAt:     timestamppb.New(pMessage.CreateAt),
	}
}
