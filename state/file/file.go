package filestate

import (
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"

	"github.com/duy-ly/nomios-go/logger"
)

type StateData struct {
	Pos string `json:"pos"`
}

type fileState struct {
	filePath string
	pos      atomic.Pointer[string]
}

func NewFileState() (*fileState, error) {
	cfg := loadConfig()

	s := new(fileState)
	s.filePath = cfg.Path
	s.pos = atomic.Pointer[string]{}

	fs, err := os.Stat(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.pos.Store(new(string))
			return s, nil
		}

		if fs.IsDir() {
			return nil, errors.New("state path is dir not file")
		}
	}

	data, err := os.ReadFile(s.filePath)
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

	s.pos.Store(&newPos)

	d, err := json.Marshal(&StateData{
		Pos: newPos,
	})
	if err != nil {
		logger.NomiosLog.Error("Error when marshal state data ", err)
		return
	}

	err = os.WriteFile(s.filePath, d, 0644)
	if err != nil {
		logger.NomiosLog.Error("Error when store state to file ", err)
		return
	}
}

func (s *fileState) GetLastPos() string {
	return *s.pos.Load()
}
