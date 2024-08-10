package services

import (
	messagePb "github.com/BigNoseCattyHome/aorb/backend/rpc/message"
	userPb "github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
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
