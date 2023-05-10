package state

type StateData struct {
	Pos string `json:"pos"`
}

type State interface {
	SaveLastPos(string)
	GetLastPos() string
}
