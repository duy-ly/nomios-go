package model

import (
	"fmt"
	"hash/fnv"
	"strings"
)

type Row struct {
	Metadata Metadata

	PrimaryKeys []string
	Values      map[string]interface{}
}

func (r *Row) GetPrimaryValues() string {
	res := ""

	for _, c := range r.PrimaryKeys {
		v := r.Values[c]

		res = fmt.Sprintf("%s:%v", res, v)
	}

	return strings.TrimPrefix(res, ":")
}

type NomiosEvent struct {
	Metadata  Metadata
	Timestamp int64

	Before Row
	After  Row
}

func (e *NomiosEvent) GetPartitionKey() string {
	return fmt.Sprintf("%s:%s:%s", e.Metadata.GetTableName(), e.Metadata.GetDatabaseName(), e.After.GetPrimaryValues())
}

func (e *NomiosEvent) GetPartitionID(size int) int {
	key := e.GetPartitionKey()

	h := fnv.New32a()
	h.Write([]byte(key))

	return int(h.Sum32()) % size
}
