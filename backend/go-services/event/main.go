package main

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/event/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/gorse"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial(rabbitmq.BuildMqConnAddr())
	exitOnError(err)

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		exitOnError(err)
	}(conn)

	tp, err := tracing.SetTraceProvider(config.Event)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	ch, err := conn.Channel()
	exitOnError(err)

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		exitOnError(err)
	}(ch)

	err = ch.ExchangeDeclare(
		strings.EventExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	exitOnError(err)

	q, err := ch.QueueDeclare(
		"event_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	exitOnError(err)

	err = ch.Qos(1, 0, false)
	exitOnError(err)

	err = ch.QueueBind(
		q.Name,
		"video.#",
		strings.EventExchange,
		false,
		nil)

	exitOnError(err)
	go Consume(ch, q.Name)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func Consume(ch *amqp.Channel, queueName string) {
	msg, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msg {
		ctx := rabbitmq.ExtractAMQPHeaders(context.Background(), d.Headers)
		ctx, span := tracing.Tracer.Start(ctx, "EventSystem")
		logger := logging.LogService("EventSystem.Recommend").WithContext(ctx)
		logging.SetSpanWithHostname(span)

		var raw models.RecommendEvent
		if err := json.Unmarshal(d.Body, &raw); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when unmarshaling the prepare json body.")
			logging.SetSpanError(span, err)
			continue
		}

		switch raw.Type {
		case 1:
			var types string
			switch raw.Source {
			case config.PollRpcServerName:
				types = "read"
			}
			var feedbacks []gorse.Feedback
			for _, id := range raw.PollId {
				feedbacks = append(feedbacks, gorse.Feedback{
					FeedbackType: types,
					UserId:       strconv.Itoa(int(raw.ActorId)),
					ItemId:       strconv.Itoa(int(id)),
					Timestamp:    time.Now().UTC().Format(time.RFC3339),
				})
			}

			if _, err := gorseClient.InsertFeedback(ctx, feedbacks); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when insert the feedback")
				logging.SetSpanError(span, err)
			}
			logger.WithFields(logrus.Fields{
				"ids": raw.PollId,
			}).Infof("Event dealt with type 1")
			span.End()
			err = d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when ack")
				logging.SetSpanError(span, err)
			}
		case 2:
			var types string
			switch raw.Source {
			case config.CommentRpcServerName:
				types = "comment"
				//case config.FavoriteRpcServerName:
				//	types = "favorite"
			}
			var feedbacks []gorse.Feedback
			for _, id := range raw.PollId {
				feedbacks = append(feedbacks, gorse.Feedback{
					FeedbackType: types,
					UserId:       strconv.Itoa(int(raw.ActorId)),
					ItemId:       strconv.Itoa(int(id)),
					Timestamp:    time.Now().UTC().Format(time.RFC3339),
				})
			}

			if _, err := gorseClient.InsertFeedback(ctx, feedbacks); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when insert the feedback")
				logging.SetSpanError(span, err)
			}
			logger.WithFields(logrus.Fields{
				"ids": raw.PollId,
			}).Infof("Event dealt with type 2")
			span.End()
			err = d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when ack")
				logging.SetSpanError(span, err)
			}
		case 3:
			var items []gorse.Item
			for _, id := range raw.PollId {
				items = append(items, gorse.Item{
					ItemId:     strconv.Itoa(int(id)),
					IsHidden:   false,
					Labels:     raw.Tag,
					Categories: raw.Category,
					Timestamp:  time.Now().UTC().Format(time.RFC3339),
					Comment:    raw.Title,
				})
			}

			if _, err := gorseClient.InsertItems(ctx, items); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when insert the items")
				logging.SetSpanError(span, err)
			}
			logger.WithFields(logrus.Fields{
				"ids":     raw.PollId,
				"tag":     raw.Tag,
				"comment": raw.Title,
			}).Infof("Event dealt with type 3")
			span.End()
			err = d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when ack")
				logging.SetSpanError(span, err)
			}
		}
	}
}

var gorseClient *gorse.GorseClient

func init() {
	gorseClient = gorse.NewGorseClient(config.Conf.Gorse.GorseAddr, config.Conf.Gorse.GorseApikey)
}
