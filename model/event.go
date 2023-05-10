package model

type Row struct {
	Metadata Metadata

	Values map[string]interface{}
}

type NomiosEvent struct {
	Metadata  Metadata
	Timestamp int64

	Before Row
	After  Row
}
