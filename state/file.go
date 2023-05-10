package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/duy-ly/nomios-go/logger"
)

type fileState struct {
	fileName string
	pos      atomic.Pointer[string]
}

func NewFileState(name string) (State, error) {
	s := new(fileState)
	s.fileName = name
	s.pos = atomic.Pointer[string]{}

	fs, err := os.Stat(s.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}

		if fs.IsDir() {
			return nil, errors.New("state path is dir not file")
		}
	}

	data, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	sd := StateData{}
	err = json.Unmarshal(data, &sd)
	if err != nil {
		return nil, err
	}

	s.pos.Store(&sd.Pos)

	return s, nil
}

func (s *fileState) SaveLastPos(newPos string) {
	if newPos == "" || *s.pos.Load() == newPos {
		return
	}

	fmt.Println(newPos)

	s.pos.Store(&newPos)

	d, err := json.Marshal(&StateData{
		Pos: newPos,
	})
	if err != nil {
		logger.NomiosLog.Error("Error when marshal state data", err)
		return
	}

	err = os.WriteFile(s.fileName, d, 0644)
	if err != nil {
		logger.NomiosLog.Error("Error when store state to file", err)
		return
	}
}

func (s *fileState) GetLastPos() string {
	return *s.pos.Load()
}
