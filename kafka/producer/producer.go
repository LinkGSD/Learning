package producer

import (
	"Learning/kafka/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

var (
	topic    = "person"
	Producer *kafka.Writer
)

func init() {
	Producer = &kafka.Writer{
		Addr:         kafka.TCP("192.168.182.132:9092"),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
	}
}

func SendMessage(ctx context.Context, user *model.User) {
	msgContent, err := json.Marshal(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("json marshal user err, user:%v,err:%v", user, err))
	}
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", user.Id)),
		Value: msgContent,
	}
	err = Producer.WriteMessages(ctx, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("写入kafka失败, user:%v,err:%v", user, err))
	}
}
