package hyperloop

import (
	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/source"
)

type HyperloopConfig struct {
	SourceConfig source.MySQLSourceConfig
	PoolConfig   consumerpool.PoolConfig

	EventStreamSize int
}
