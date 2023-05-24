package publisher

import (
	"github.com/duy-ly/nomios-go/model"
	kafkapublisher "github.com/duy-ly/nomios-go/publisher/kafka"
	"github.com/spf13/viper"
)

type Publisher interface {
	Publish(msg []*model.NomiosEvent) error
	Close() error
}

func NewPublisher() (Publisher, error) {
	kind := viper.GetString("publisher.kind")
	if kind == "" {
		kind = "kafka"
	}

	switch kind {
	default:
		return kafkapublisher.NewKafkaPublisher()
	}
}
