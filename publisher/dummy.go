package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/duy-ly/nomios-go/model"
)

type dummyPublisher struct {
}

func NewDummyPublisher() Publisher {
	p := new(dummyPublisher)

	return p
}

func (p *dummyPublisher) Publish(msg []*model.NomiosEvent) error {
	// TODO: implement publish
	m, _ := json.Marshal(msg)
	fmt.Println("events", string(m))

	return nil
}

func (p *dummyPublisher) Close() error {
	return nil
}
