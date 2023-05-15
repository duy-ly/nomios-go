package hyperloop

import (
	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/source"
	"github.com/duy-ly/nomios-go/state"
)

type Hyperloop struct {
	running bool

	cp consumerpool.ConsumerPool
	sm state.StateManager
	s  source.Source
}

func NewHyperloop() *Hyperloop {
	stt, err := state.NewState()
	if err != nil {
		logger.NomiosLog.Panic("Error when init state ", err)
	}
	cp, err := consumerpool.NewConsumerPool()
	if err != nil {
		logger.NomiosLog.Panic("Error when init consumer pool ", err)
	}
	sm := state.NewStateManager(stt, cp.GetLastGTID)
	s, err := source.NewSource()
	if err != nil {
		logger.NomiosLog.Panic("Error when init source ", err)
	}

	h := new(Hyperloop)
	h.cp = cp
	h.sm = sm
	h.s = s

	return h
}

func (h *Hyperloop) Start() {
	h.sm.Start()

	h.cp.Start()

	h.s.Start(h.sm.GetLastCheckpoint(), h.cp.GetStream())

	h.running = true
}

func (h *Hyperloop) Stop() {
	h.sm.Stop()

	// stop source first to prevent push more event to stream
	h.s.Stop()

	h.cp.Stop()

	h.sm.Checkpoint()
}

func (h *Hyperloop) IsRunning() bool {
	return h.running
}

func (h *Hyperloop) OnError(err error) {

}
