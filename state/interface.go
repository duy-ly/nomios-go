package state

import (
	filestate "github.com/duy-ly/nomios-go/state/file"
)

type State interface {
	SaveLastPos(string)
	GetLastPos() string
}

func NewState(kind string) (State, error) {
	switch kind {
	default:
		return filestate.NewFileState()
	}
}
