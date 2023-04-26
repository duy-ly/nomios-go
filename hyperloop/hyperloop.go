package hyperloop

import (
	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/event"
	"github.com/duy-ly/nomios-go/source"
)

type Hyperloop struct {
	stream  chan event.NomiosEvent
	cfg     HyperloopConfig
	running bool

	cp *consumerpool.ConsumerPool
	s  *source.Source
}

func NewHyperloop(cfg HyperloopConfig) *Hyperloop {
	if cfg.EventStreamSize == 0 {
		cfg.EventStreamSize = 20000
	}

	h := new(Hyperloop)
	h.cfg = cfg
	h.stream = make(chan event.NomiosEvent, cfg.EventStreamSize)

	h.cp = consumerpool.NewConsumerPool(cfg.PoolConfig)
	s, err := source.NewSource(h.cfg.SourceConfig)
	if err != nil {
		panic(err)
	}
	h.s = s

	return h
}

func (h *Hyperloop) Start() {
	// TODO: schedule state manager

	h.cp.Start(h.stream)

	h.s.Start("", 0, h.stream)

	h.running = true
}

func (h *Hyperloop) Stop() {
	// TODO: stop state manager

	// stop source first to prevent push more event to stream
	h.s.Stop()

	close(h.stream)

	h.cp.Stop()

	// TODO: manual state checkpoint
}

func (h *Hyperloop) IsRunning() bool {
	return h.running
}

func (h *Hyperloop) OnError(err error) {

}
