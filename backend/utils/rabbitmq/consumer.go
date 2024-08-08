package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeMessage(ctx context.Context, queueName string) (msg <-chan amqp.Delivery, err error) {
	ch, err := RabbitMQ.Channel()
	if err != nil {
		return nil, err
	}
	q, _ := ch.QueueDeclare(queueName, true, false, false, false, nil)
	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}
	return ch.Consume(q.Name, "", false, false, false, false, nil)
}
