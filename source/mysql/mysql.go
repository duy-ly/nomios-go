package mysqlsource

import (
	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
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
	syncerCfg.IncludeTableRegex = buildIncludeTableRegex(cfg.Database, cfg.TableIncludeList)
	syncerCfg.ExcludeTableRegex = []string{"mysql\\..*"}
	syncerCfg.ParseTime = true
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

func (s *mysqlSource) Start(lastID string, stream chan []*model.NomiosEvent) {
	go func() {
		s.syncer.SetEventHandler(NewEventMapperHandler(stream))

		// Start canal
		gtidSet, err := mysql.ParseGTIDSet("mysql", lastID)
		if err != nil {
			logger.NomiosLog.Panic("Error when parse gtid set ", err)
		}

		err = s.syncer.StartFromGTID(gtidSet)
		if err != nil {
			logger.NomiosLog.Panic("Error when start from gtid ", err)
		}

		<-s.stopSig

		s.syncer.Close()
	}()
}

func (s *mysqlSource) Stop() {
	s.stopSig <- true
}
