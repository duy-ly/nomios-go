package consumerpool

import "github.com/duy-ly/nomios-go/event"

type Publisher struct {
}

func NewPublisher() *Publisher {
	p := new(Publisher)

	return p
}

func (p *Publisher) Publish(msg []event.NomiosEvent) error {
	// TODO

	return nil
}
