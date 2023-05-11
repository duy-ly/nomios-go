package mysqlsource

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

	gtid   atomic.Pointer[string]
	stream chan []*model.NomiosEvent
}

func NewEventMapperHandler(stream chan []*model.NomiosEvent) canal.EventHandler {
	return &EventMapperHandler{
		gtid:   atomic.Pointer[string]{},
		stream: stream,
	}
}

func (h *EventMapperHandler) OnRotate(header *replication.EventHeader, e *replication.RotateEvent) error {
	logger.NomiosLog.Infof("Rotate event: - Position: %d;- NextName: %s", e.Position, e.NextLogName)

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
			Pos: e.Header.LogPos,
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

		row := h.generateRow(r, e.Table)

		switch e.Action {
		case "insert":
			ne := base
			ne.After = *row

			se = append(se, &ne)
		case "update":
			after := h.generateRow(e.Rows[i+1], e.Table)

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

func (h *EventMapperHandler) generateRow(values []interface{}, table *schema.Table) *model.Row {
	primaryKeys := make([]string, 0)
	rowData := make(map[string]interface{})

	if len(values) != len(table.Columns) {
		logger.NomiosLog.Errorf("Row length %d and column length %d don't match", len(values), len(table.Columns))
		return nil
	}

	for _, i := range table.Indexes {
		if i.Name == "PRIMARY" {
			primaryKeys = i.Columns
			break
		}
	}

	for i, v := range values {
		c := table.Columns[i]

		rowData[c.Name] = v
	}

	return &model.Row{
		PrimaryKeys: primaryKeys,
		Values:      rowData,
	}
}

func (h *EventMapperHandler) OnGTID(header *replication.EventHeader, gtidSet mysql.GTIDSet) error {
	newGtid := gtidSet.String()
	logger.NomiosLog.Info("Start new gtid ", newGtid)

	h.gtid.Store(&newGtid)

	return nil
}

func (h *EventMapperHandler) String() string {
	return "EventMapperHandler"
}
