package source

import (
	"sync/atomic"
	"time"

	"github.com/duy-ly/nomios-go/logger"
	"github.com/duy-ly/nomios-go/model"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
)

type EventMapperHandler struct {
	canal.DummyEventHandler

	gtid    atomic.Pointer[string]
	posName atomic.Pointer[string]
	stream  chan []*model.NomiosEvent
}

func NewEventMapperHandler(stream chan []*model.NomiosEvent) canal.EventHandler {
	return &EventMapperHandler{
		gtid:    atomic.Pointer[string]{},
		posName: atomic.Pointer[string]{},
		stream:  stream,
	}
}

func (h *EventMapperHandler) OnRotate(header *replication.EventHeader, e *replication.RotateEvent) error {
	logger.NomiosLog.Infof("Rotate event: - Position: %d;- NextName: %s", e.Position, e.NextLogName)

	logName := string(e.NextLogName)
	h.posName.Store(&logName)

	return nil
}

func (h *EventMapperHandler) OnRow(e *canal.RowsEvent) error {
	eventMetadata := model.MySQLMetadata{
		GTID:     *h.gtid.Load(),
		ServerID: int64(e.Header.ServerID),
		Ts:       int64(e.Header.Timestamp),
		TableSchema: model.MySQLTableSchema{
			Name:     e.Table.Name,
			Database: e.Table.Schema,
		},
		BinlogPosition: mysql.Position{
			Name: *h.posName.Load(),
			Pos:  e.Header.LogPos,
		},
	}

	se := make([]*model.NomiosEvent, 0)

	base := model.NomiosEvent{
		Metadata:  &eventMetadata,
		Timestamp: time.Now().UnixMilli(),
	}

	for i, r := range e.Rows {
		if i%2 == 1 {
			// not process even event due to row change are tuple [before, after]
			continue
		}

		row := h.generateRow(r, e.Table.Columns)

		switch e.Action {
		case "insert":
			ne := base
			ne.After = *row

			se = append(se, &ne)
		case "update":
			after := h.generateRow(e.Rows[i+1], e.Table.Columns)

			ne := base
			ne.Before = *row
			ne.After = *after

			se = append(se, &ne)
		case "delete":
			ne := base
			ne.Before = *row

			se = append(se, &ne)
		}

	}

	h.stream <- se

	return nil
}

func (h *EventMapperHandler) generateRow(values []interface{}, cols []schema.TableColumn) *model.Row {
	rowData := make(map[string]interface{})

	if len(values) != len(cols) {
		logger.NomiosLog.Errorf("Row length %d and column length %d don't match", len(values), len(cols))
		return nil
	}

	for i, v := range values {
		c := cols[i]

		rowData[c.Name] = v
	}

	return &model.Row{
		Values: rowData,
	}
}

func (h *EventMapperHandler) OnXID(header *replication.EventHeader, pos mysql.Position) error {
	return nil
}

func (h *EventMapperHandler) OnGTID(header *replication.EventHeader, gtidSet mysql.GTIDSet) error {
	newGtid := gtidSet.String()
	logger.NomiosLog.Info("Start new gtid", newGtid)

	h.gtid.Store(&newGtid)

	return nil
}

func (h *EventMapperHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, gtidSet mysql.GTIDSet, isForce bool) error {
	return nil
}

func (h *EventMapperHandler) String() string {
	return "EventMapperHandler"
}
