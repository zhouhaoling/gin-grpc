package kk

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type LogData struct {
	Topic string
	//接收json数据
	Data []byte
}

type KafkaWriter struct {
	w    *kafka.Writer
	data chan LogData
}

func GetWriter(addr string) *KafkaWriter {
	w := &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Balancer: &kafka.LeastBytes{},
	}
	k := &KafkaWriter{
		w:    w,
		data: make(chan LogData, 100),
	}
	go k.sendKafka()
	return k
}

func (kw *KafkaWriter) Send(data LogData) {
	kw.data <- data
}

func (kw *KafkaWriter) Close() {
	if kw.w != nil {
		kw.w.Close()
	}
}

func (kw *KafkaWriter) sendKafka() {
	for {
		select {
		case data := <-kw.data:
			message := []kafka.Message{
				{
					Topic: data.Topic,
					Key:   []byte("logMsg"),
					Value: data.Data,
				},
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			for i := 0; i < retries; i-- {

				err = kw.w.WriteMessages(ctx, message...)
				if err == nil {
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}

				if err != nil {
					log.Printf("kafka send writemessage err %s \n", err)
				}
			}

		}
	}
}
