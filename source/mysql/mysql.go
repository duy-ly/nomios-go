package mysqlsource

import (
	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/util"
	"github.com/go-mysql-org/go-mysql/canal"
)

type mysqlSource struct {
	stopSig chan bool
	syncer  *canal.Canal
}

// NewMySQLSource -- create new MySQL source from config
func NewMySQLSource() (*mysqlSource, error) {
	cfg := loadConfig()

	syncerCfg := canal.NewDefaultConfig()
	syncerCfg.Addr = cfg.Addr
	syncerCfg.Flavor = "mysql"
	syncerCfg.User = cfg.User
	syncerCfg.Password = cfg.Pass
	syncerCfg.Dump.ExecutionPath = ""
	syncerCfg.Logger = logger.NomiosLog
	syncerCfg.ServerID = cfg.ServerID

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

func (s *mysqlSource) Start(posStr string, stream chan []*model.NomiosEvent) {
	go func() {
		pos := util.ParseEventPos(posStr)

		s.syncer.SetEventHandler(NewEventMapperHandler(pos.Name, stream))

		// Start canal
		err := s.syncer.RunFrom(pos)
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
