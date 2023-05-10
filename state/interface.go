package state

import (
	filestate "github.com/duy-ly/nomios-go/state/file"
	"github.com/spf13/viper"
)

type State interface {
	SaveLastID(string)
	GetLastID() string
}

func NewState() (State, error) {
	kind := viper.GetString("state.kind")
	if kind == "" {
		kind = "file"
	}

	switch kind {
	default:
		return filestate.NewFileState()
	}
}
