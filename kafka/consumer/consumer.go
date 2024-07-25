package consumer

import (
	"Learning/kafka/model"
	"Learning/kafka/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

var (
	topic     = "user"
	Consumer  *kafka.Reader
	DeadQueue *kafka.Writer
)

func init() {
	Consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:                []string{"192.168.182.132:9092"},
		Topic:                  topic,
		GroupID:                "t22",
		GroupTopics:            nil,
		Partition:              0,
		Dialer:                 nil,
		QueueCapacity:          0,
		MinBytes:               0, // 10KB
		MaxBytes:               0, // 10MB
		ReadBatchTimeout:       0,
		ReadLagInterval:        0,
		GroupBalancers:         nil,
		HeartbeatInterval:      0,
		CommitInterval:         time.Second,
		PartitionWatchInterval: 0,
		WatchPartitionChanges:  false,
		SessionTimeout:         0,
		RebalanceTimeout:       0,
		JoinGroupBackoff:       0,
		RetentionTime:          0,
		StartOffset:            kafka.FirstOffset,
		ReadBackoffMax:         0,
		ReadBackoffMin:         0,
		Logger:                 nil,
		ErrorLogger:            nil,
		IsolationLevel:         0,
		MaxAttempts:            0,
		OffsetOutOfRangeError:  false,
	})

	DeadQueue = &kafka.Writer{
		Addr:         kafka.TCP("192.168.182.132:9092"),
		Topic:        fmt.Sprintf("dead-%s", topic),
		RequiredAcks: kafka.RequireOne,
	}
}

func ReadMessage(ctx context.Context) {
	for {
		if msg, err := Consumer.ReadMessage(ctx); err != nil {
			fmt.Println(fmt.Sprintf("读取kafka失败,err:%v", err))
		} else {
			user := &model.User{}
			err := json.Unmarshal(msg.Value, user)
			if err != nil {
				fmt.Println(fmt.Sprintf("json unmarshal msg value err, msg:%v,err:%v", user, err))
			}
			//幂等
			if redis.IsConsumed(fmt.Sprintf("%s-%d", topic, user.Id)) {
				fmt.Println(fmt.Sprintf("topic=%s,partition=%d,offset=%d,key=%s,user=%v", msg.Topic, msg.Partition, msg.Offset, msg.Key, user))
				if err != nil {
				retry:
					err := DeadQueue.WriteMessages(ctx, kafka.Message{
						Key:     msg.Key,
						Value:   msg.Value,
						Headers: []kafka.Header{{"errors", []byte(err.Error())}},
					})
					if err != nil {
						goto retry
					}
				}
			} else {
				fmt.Println(fmt.Sprintf("user:%v 已经消费过,不再消费", user))
			}
		}
	}
}
