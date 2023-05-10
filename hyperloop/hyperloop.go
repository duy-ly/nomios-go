package hyperloop

import (
	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/source"
	"github.com/duy-ly/nomios-go/state"
)

type Hyperloop struct {
	stream  chan []*model.NomiosEvent
	running bool

	cp *consumerpool.ConsumerPool
	sm *state.StateManager
	s  source.Source
}

func NewHyperloop() *Hyperloop {
	cfg := loadConfig()

	stt, err := state.NewState(cfg.StateKind)
	if err != nil {
		logger.NomiosLog.Panic("Error when init state ", err)
	}
	cp, err := consumerpool.NewConsumerPool(cfg.PublisherKind)
	if err != nil {
		logger.NomiosLog.Panic("Error when init consumer pool ", err)
	}
	s, err := source.NewSource(cfg.SourceKind)
	if err != nil {
		logger.NomiosLog.Panic("Error when init source ", err)
	}

	h := new(Hyperloop)
	h.stream = make(chan []*model.NomiosEvent, cfg.EventStreamSize)
	h.cp = cp
	h.sm = state.NewStateManager(stt, cp.GetLastEventPos)
	h.s = s

	return h
}

func (h *Hyperloop) Start() {
	h.sm.Start()

	h.cp.Start(h.stream)

	pos := h.sm.GetState().GetLastPos()
	h.s.Start(pos, h.stream)

	h.running = true
}

func (h *Hyperloop) Stop() {
	h.sm.Stop()

	// stop source first to prevent push more event to stream
	h.s.Stop()

	close(h.stream)

	h.cp.Stop()

	h.sm.Checkpoint()
}

func (h *Hyperloop) IsRunning() bool {
	return h.running
}

func (h *Hyperloop) OnError(err error) {

}
