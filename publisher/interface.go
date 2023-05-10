package publisher

import "github.com/duy-ly/nomios-go/model"

type Publisher interface {
	Publish(msg []*model.NomiosEvent) error
	Close() error
}
