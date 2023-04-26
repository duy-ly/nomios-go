package source

import (
	"github.com/duy-ly/nomios-go/event"
	"github.com/go-mysql-org/go-mysql/replication"
)

type SourceEventMapper interface {
	MapEvent(*replication.BinlogEvent) []event.NomiosEvent
}

type sourceEventMapper struct {
}

func NewSourceEventMapper() SourceEventMapper {
	m := new(sourceEventMapper)

	return m
}

func (m *sourceEventMapper) MapEvent(binlogEvent *replication.BinlogEvent) []event.NomiosEvent {
	return make([]event.NomiosEvent, 0)
}
