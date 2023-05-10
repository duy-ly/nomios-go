package source

import (
	"github.com/duy-ly/nomios-go/model"
	mysqlsource "github.com/duy-ly/nomios-go/source/mysql"
)

type Source interface {
	Start(pos string, stream chan []*model.NomiosEvent)
	Stop()
}

func NewSource(kind string) (Source, error) {
	switch kind {
	default:
		return mysqlsource.NewMySQLSource()
	}
}
