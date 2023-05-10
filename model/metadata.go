package model

import "github.com/go-mysql-org/go-mysql/mysql"

type Metadata interface {
	GetPos() mysql.Position
}

type MySQLMetadata struct {
	GTID           string
	ServerID       int64
	Ts             int64
	TableSchema    MySQLTableSchema
	BinlogPosition mysql.Position
}

func (m *MySQLMetadata) GetPos() mysql.Position {
	return m.BinlogPosition
}

type MySQLTableSchema struct {
	Name     string
	Database string
}
