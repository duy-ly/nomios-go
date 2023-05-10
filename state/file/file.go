package filestate

import (
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"

	"github.com/duy-ly/nomios-go/logger"
)

type StateData struct {
	ID string `json:"id"`
}

type fileState struct {
	filePath string
	lastID   atomic.Pointer[string]
}

func NewFileState() (*fileState, error) {
	cfg := loadConfig()

	s := new(fileState)
	s.filePath = cfg.Path
	s.lastID = atomic.Pointer[string]{}

	fs, err := os.Stat(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.lastID.Store(new(string))
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

	s.lastID.Store(&sd.ID)

	return s, nil
}

func (s *fileState) SaveLastID(newID string) {
	if newID == "" || *s.lastID.Load() == newID {
		return
	}

	s.lastID.Store(&newID)

	d, err := json.Marshal(&StateData{
		ID: newID,
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

func (s *fileState) GetLastID() string {
	return *s.lastID.Load()
}
