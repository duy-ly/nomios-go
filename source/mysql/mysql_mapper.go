package mysqlsource

import (
	"fmt"
	"strconv"
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
	primaryKeys := make([]string, 0)

	for _, c := range e.Table.PKColumns {
		primaryKeys = append(primaryKeys, e.Table.Columns[c].Name)
	}

	eventMetadata := model.MySQLMetadata{
		GTID:     *h.gtid.Load(),
		ServerID: int64(e.Header.ServerID),
		Ts:       int64(e.Header.Timestamp),
		TableSchema: model.MySQLTableSchema{
			Database:    e.Table.Schema,
			Name:        e.Table.Name,
			PrimaryKeys: primaryKeys,
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

	if len(se) > 0 {
		h.stream <- se
	}

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

		switch c.Type {
		case schema.TYPE_NUMBER, schema.TYPE_MEDIUM_INT:
			switch tv := v.(type) {
			case int:
				v = int64(tv)
			case int8:
				v = int64(tv)
			case int16:
				v = int64(tv)
			case int32:
				v = int64(tv)
			case uint:
				v = int64(tv)
			case uint8:
				v = int64(tv)
			case uint16:
				v = int64(tv)
			case uint32:
				v = int64(tv)
			case uint64:
				v = int64(tv)
			default:
			}
		case schema.TYPE_ENUM:
			vi, _ := v.(int64)
			if len(c.EnumValues) < int(vi-1) {
				continue
			}

			v = c.EnumValues[vi-1]
		case schema.TYPE_DATETIME:
		// 	vs, _ := v.(string)
		// 	newV, err := time.ParseInLocation(time.DateTime, vs, time.UTC)
		// 	if err == nil {
		// 		v = newV
		// 	}
		case schema.TYPE_TIMESTAMP:
		// 	vs, _ := v.(string)
		// 	newV, err := time.ParseInLocation(time.DateTime, vs, time.Local)
		// 	if err == nil {
		// 		v = newV
		// 	}
		case schema.TYPE_DATE:
			vs, _ := v.(string)
			newV, err := time.Parse(time.DateOnly, vs)
			if err == nil {
				v = newV
			}
		case schema.TYPE_TIME:
			vs, _ := v.(string)
			newV, err := time.Parse(time.TimeOnly, vs)
			if err == nil {
				v = newV.AddDate(1970, 0, 0)
			}
		case schema.TYPE_JSON:

		case schema.TYPE_DECIMAL:
			if v != nil {
				vs, _ := v.(string)
				newV, err := strconv.ParseFloat(vs, 64)
				if err == nil {
					v = newV
				}
			}
		default:
			fmt.Println(c)
		}

		rowData[c.Name] = v
	}

	return &model.Row{
		Values: rowData,
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
