package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/poll/services"
	pollPb "github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/messageQueue"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
)

type SyncPoll struct {
}

func (p *SyncPoll) SyncPoll2Redis(ctx context.Context, queueName string) error {
	msg, err := rabbitmq.ConsumeMessage(ctx, queueName)
	if err != nil {
		return err
	}
	// 用来阻塞这个接口，使其永远不会返回
	// forever里面不能传任何东西，否则阻塞失效
	var forever chan struct{}
	go func() {
		for d := range msg {
			// 落库
			var rPoll *pollPb.Poll
			err = json.Unmarshal(d.Body, &rPoll)
			fmt.Println(rPoll)
			if err != nil {
				return
			}
			err = services.PollMQ2Redis(ctx, rPoll)
			if err != nil {
				return
			}
			d.Ack(false)
		}
	}()
	<-forever // 阻塞
	return nil
}

func SyncPoll2Redis(ctx context.Context) {
	Sync := new(SyncPoll)
	err := Sync.SyncPoll2Redis(ctx, messageQueue.Poll2RedisQueue)
	if err != nil {
		return
	}
}
