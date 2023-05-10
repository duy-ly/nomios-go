package state

import (
	"time"
)

type StateManager struct {
	s State

	stopSig   chan bool
	collector func() string
}

func NewStateManager(s State, fn func() string) *StateManager {
	m := new(StateManager)
	m.s = s
	m.stopSig = make(chan bool, 1)
	m.collector = fn

	return m
}

func (m *StateManager) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-m.stopSig:
				return
			case <-ticker.C:
				m.Checkpoint()
			default:
			}
		}
	}()
}

func (m *StateManager) Checkpoint() {
	if m.collector == nil {
		return
	}

	m.s.SaveLastPos(m.collector())
}

func (m *StateManager) GetState() State {
	return m.s
}

func (m *StateManager) Stop() {
	m.stopSig <- true
}
