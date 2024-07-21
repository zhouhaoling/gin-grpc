package kk

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	r *kafka.Reader
}

func (kr *KafkaReader) readMsg() {
	for {
		message, err := kr.r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("kafka readMsg err %s \n", err.Error())
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
	}
}

func GetReader(brokers []string, groupId, topic string) *KafkaReader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId, //同一个组下的consumer 协同工作 共同消费topic队列里的内容
		Topic:    topic,
		MinBytes: 10e3, //10KB
		MaxBytes: 10e6, //10MB
	})
	k := &KafkaReader{
		r: r,
	}
	go k.readMsg()
	return k
}
