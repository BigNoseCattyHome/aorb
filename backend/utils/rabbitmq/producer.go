package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func SendMessage2MQ(body []byte, queueName string) (err error) {
	ch, err := RabbitMQ.Channel()
	if err != nil {
		return
	}
	q, _ := ch.QueueDeclare(queueName, true, false, false, false, nil)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		Body:         body,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return
	}
	return
}
