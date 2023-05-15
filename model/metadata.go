package model

import "github.com/go-mysql-org/go-mysql/mysql"

type Metadata interface {
	GetID() string
	GetTableName() string
	GetDatabaseName() string
	GetPrimaryKeys() []string
	GetPos() mysql.Position
}

type MySQLMetadata struct {
	GTID           string           `json:"gtid"`
	ServerID       int64            `json:"server_id"`
	Ts             int64            `json:"ts"`
	TableSchema    MySQLTableSchema `json:"schema"`
	BinlogPosition mysql.Position   `json:"-"`
}

type MySQLTableSchema struct {
	Database    string   `json:"database"`
	Name        string   `json:"name"`
	PrimaryKeys []string `json:"primary_keys"`
}

func (m *MySQLMetadata) GetID() string {
	return m.GTID
}

func (m *MySQLMetadata) GetDatabaseName() string {
	return m.TableSchema.Database
}

func (m *MySQLMetadata) GetTableName() string {
	return m.TableSchema.Name
}

func (m *MySQLMetadata) GetPrimaryKeys() []string {
	return m.TableSchema.PrimaryKeys
}

func (m *MySQLMetadata) GetPos() mysql.Position {
	return m.BinlogPosition
}
