package source

import (
	"github.com/duy-ly/nomios-go/event"
	"github.com/go-mysql-org/go-mysql/replication"
)

type SourceEventMapper interface {
	MapEvent(*replication.BinlogEvent) []event.NomiosEvent
}
