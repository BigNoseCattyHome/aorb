package rabbitmq

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQ *amqp.Connection

func BuildMqConnAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.Conf.RabbitMQ.Username,
		config.Conf.RabbitMQ.Password,
		config.Conf.RabbitMQ.Host,
		config.Conf.RabbitMQ.Port,
		config.Conf.RabbitMQ.VhostPrefix)
}

func init() {
	mqAddr := BuildMqConnAddr()
	conn, err := amqp.Dial(mqAddr)
	if err != nil {
		panic(err)
	}
	RabbitMQ = conn
}
