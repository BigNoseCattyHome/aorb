package rabbitmq

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
)

func BuildMqConnAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.Conf.RabbitMQ.Username,
		config.Conf.RabbitMQ.Password,
		config.Conf.RabbitMQ.Host,
		config.Conf.RabbitMQ.Port,
		config.Conf.RabbitMQ.VhostPrefix)
}
