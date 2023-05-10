package source

import (
	"fmt"

	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/util"
	"github.com/go-mysql-org/go-mysql/canal"
)

type MySQLSourceConfig struct {
	// mysql binlog config
	Host string
	Port uint16
	User string
	Pass string

	// syncer config
	ServerID uint32
}

type mysqlSource struct {
	stopSig chan bool
	syncer  *canal.Canal
}

// NewMySQLSource -- create new MySQL source from config
func NewMySQLSource(cfg MySQLSourceConfig) (Source, error) {
	syncerCfg := canal.NewDefaultConfig()
	syncerCfg.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	syncerCfg.Flavor = "mysql"
	syncerCfg.User = cfg.User
	syncerCfg.Password = cfg.Pass
	syncerCfg.Dump.ExecutionPath = ""
	syncerCfg.Logger = logger.NomiosLog

	// create instance
	s := new(mysqlSource)
	s.stopSig = make(chan bool, 1)

	var err error
	s.syncer, err = canal.NewCanal(syncerCfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *mysqlSource) Start(pos string, stream chan []*model.NomiosEvent) {
	go func() {
		s.syncer.SetEventHandler(NewEventMapperHandler(stream))

		// Start canal
		err := s.syncer.RunFrom(util.ParseEventPos(pos))
		if err != nil {
			panic(err)
		}

		<-s.stopSig

		s.syncer.Close()
	}()
}

func (s *mysqlSource) Stop() {
	s.stopSig <- true
}
