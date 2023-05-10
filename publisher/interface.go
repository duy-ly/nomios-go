package publisher

import (
	"github.com/duy-ly/nomios-go/model"
	kafkapublisher "github.com/duy-ly/nomios-go/publisher/kafka"
)

type Publisher interface {
	Publish(msg []*model.NomiosEvent) error
	Close() error
}

func NewPublisher(kind string) (Publisher, error) {
	switch kind {
	default:
		return kafkapublisher.NewKafkaPublisher()
	}
}
