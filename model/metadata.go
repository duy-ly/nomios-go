package model

import "github.com/go-mysql-org/go-mysql/mysql"

type Metadata interface {
	GetID() string
	GetTableName() string
	GetDatabaseName() string
	GetPos() mysql.Position
}

type MySQLMetadata struct {
	GTID           string
	ServerID       int64
	Ts             int64
	TableSchema    MySQLTableSchema
	BinlogPosition mysql.Position
}

func (m *MySQLMetadata) GetID() string {
	return m.GTID
}

func (m *MySQLMetadata) GetTableName() string {
	return m.TableSchema.Name
}

func (m *MySQLMetadata) GetDatabaseName() string {
	return m.TableSchema.Database
}

func (m *MySQLMetadata) GetPos() mysql.Position {
	return m.BinlogPosition
}

type MySQLTableSchema struct {
	Name     string
	Database string
}
