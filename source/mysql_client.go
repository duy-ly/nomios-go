package source

import (
	"context"

	"github.com/duy-ly/nomios-go/event"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type SourceConfig struct {
	// mysql binlog config
	Host string
	Port uint16
	User string
	Pass string

	// syncer config
	ServerID uint32
}

type Source struct {
	stopSig chan bool

	mapper SourceEventMapper

	syncer *replication.BinlogSyncer
}

// NewSource -- create new source from config
func NewSource(cfg SourceConfig) (*Source, error) {
	bsCfg := replication.BinlogSyncerConfig{
		ServerID: cfg.ServerID,
		Flavor:   "mysql",
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Pass,
	}

	// create instance
	s := new(Source)
	s.stopSig = make(chan bool, 1)
	s.syncer = replication.NewBinlogSyncer(bsCfg)

	return s, nil
}

// Start -- start sync from mysql then mapping BinlogEvent to NomiosEvent
func (s *Source) Start(binlogFile string, binlogPos uint32, stream chan event.NomiosEvent) {
	go func() {
		binlogStreamer, _ := s.syncer.StartSync(mysql.Position{
			Name: binlogFile,
			Pos:  binlogPos,
		})

		for {
			select {
			case <-s.stopSig:
				break
			default:
				binlogEvent, err := binlogStreamer.GetEvent(context.Background())
				if err != nil {
					// TODO: handling error
					return
				}

				mappedEvents := s.mapper.MapEvent(binlogEvent)
				for _, e := range mappedEvents {
					stream <- e
				}
			}
		}
	}()
}
