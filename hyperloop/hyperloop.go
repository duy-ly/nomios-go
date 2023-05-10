package hyperloop

import (
	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/source"
	"github.com/duy-ly/nomios-go/state"
)

type Hyperloop struct {
	stream  chan []*model.NomiosEvent
	cfg     HyperloopConfig
	running bool

	cp *consumerpool.ConsumerPool
	sm *state.StateManager
	s  source.Source
}

func NewHyperloop(cfg HyperloopConfig) *Hyperloop {
	if cfg.EventStreamSize == 0 {
		cfg.EventStreamSize = 20000
	}

	h := new(Hyperloop)
	h.cfg = cfg
	h.stream = make(chan []*model.NomiosEvent, cfg.EventStreamSize)

	stt, err := state.NewFileState("./checkpoint.nom")
	if err != nil {
		panic(err)
	}
	h.cp = consumerpool.NewConsumerPool(cfg.PoolConfig)
	h.sm = state.NewStateManager(stt, h.cp.GetLastEventPos)
	s, err := source.NewMySQLSource(h.cfg.SourceConfig)
	if err != nil {
		panic(err)
	}
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
