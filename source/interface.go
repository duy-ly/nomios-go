package source

import "github.com/duy-ly/nomios-go/model"

type Source interface {
	Start(pos string, stream chan []*model.NomiosEvent)
	Stop()
}
