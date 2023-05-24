package state

import (
	"time"
)

type StateManager interface {
	Start()
	Checkpoint()
	GetLastCheckpoint() string
	Stop()
}

type stateManager struct {
	s State

	cron      time.Duration
	stopSig   chan bool
	collector func() string
}

func NewStateManager(s State, collectFn func() string) StateManager {
	cfg := loadConfig()

	m := new(stateManager)
	m.s = s
	m.cron = cfg.Cron
	m.stopSig = make(chan bool, 1)
	m.collector = collectFn

	return m
}

func (m *stateManager) Start() {
	go func() {
		ticker := time.NewTicker(m.cron)

		for {
			select {
			case <-m.stopSig:
				ticker.Stop()
				return
			case <-ticker.C:
				m.Checkpoint()
			default:
			}
		}
	}()
}

func (m *stateManager) Checkpoint() {
	if m.collector == nil {
		return
	}

	m.s.SaveLastID(m.collector())
}

func (m *stateManager) GetLastCheckpoint() string {
	return m.s.GetLastID()
}

func (m *stateManager) Stop() {
	m.stopSig <- true
}
