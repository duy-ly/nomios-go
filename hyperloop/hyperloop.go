package hyperloop

import (
	"github.com/duy-ly/nomios-go/event"
)

type Hyperloop struct {
	stream chan event.NomiosEvent

	running bool
}

func NewHyperloop() *Hyperloop {
	h := new(Hyperloop)
	h.stream = make(chan event.NomiosEvent, 20000)

	return h
}

func (h *Hyperloop) Start() {

}

func (h *Hyperloop) Stop() {

}

func (h *Hyperloop) IsRunning() bool {
	return h.running
}

func (h *Hyperloop) OnError() {

}
