package kafkapublisher

import (
	"context"
	"encoding/json"

	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/segmentio/kafka-go"
)

type kafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher() (*kafkaPublisher, error) {
	cfg := loadConfig()

	p := new(kafkaPublisher)
	p.writer = &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Addrs...),
		Topic:                  cfg.Topic,
		AllowAutoTopicCreation: true,
		// Balancer: &kafka.Hash{},
	}

	return p, nil
}

func (p *kafkaPublisher) Publish(msg []*model.NomiosEvent) error {
	kmsg := make([]kafka.Message, 0)

	for _, e := range msg {
		m, err := json.Marshal(e)
		if err != nil {
			logger.NomiosLog.Error("Error when marshal nomios event ", err)
			continue
		}

		kmsg = append(kmsg, kafka.Message{
			Value: m,
		})
	}

	err := p.writer.WriteMessages(context.Background(), kmsg...)
	if err != nil {
		logger.NomiosLog.Error("Error when publish message to kafka ", err)
		return err
	}

	return nil
}

func (p *kafkaPublisher) Close() error {
	return p.writer.Close()
}
