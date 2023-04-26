package event

type Metadata interface {
}

type MySQLMetadata struct {
	tableSchema    struct{}
	binlogPosition struct{}
	gtid           string
	serverID       int64
	ts             int64
}

type Row struct {
}

type NomiosEvent struct {
	metadata Metadata

	before Row
	after  Row
}
