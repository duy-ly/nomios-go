package source

import (
	"github.com/duy-ly/nomios-go/model"
	mysqlsource "github.com/duy-ly/nomios-go/source/mysql"
	"github.com/spf13/viper"
)

type Source interface {
	Start(pos string, stream chan []*model.NomiosEvent)
	Stop()
}

func NewSource() (Source, error) {
	kind := viper.GetString("source.kind")
	if kind == "" {
		kind = "mysql"
	}

	switch kind {
	default:
		return mysqlsource.NewMySQLSource()
	}
}
