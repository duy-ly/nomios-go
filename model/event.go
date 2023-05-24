package model

import (
	"fmt"
	"hash/fnv"
	"strings"
)

type Row struct {
	Metadata Metadata `json:"-"`

	Values map[string]interface{} `json:"values"`
}

func (r *Row) GetPrimaryValues(columns []string) string {
	res := ""

	for _, c := range columns {
		v := r.Values[c]

		res = fmt.Sprintf("%s:%v", res, v)
	}

	return strings.TrimPrefix(res, ":")
}

type NomiosEvent struct {
	Metadata  Metadata `json:"metadata"`
	Timestamp int64    `json:"timestamp"`

	Before Row `json:"before,omitempty"`
	After  Row `json:"after,omitempty"`
}

func (e *NomiosEvent) GetPartitionKey() string {
	return fmt.Sprintf("%s:%s:%s", e.Metadata.GetTableName(), e.Metadata.GetDatabaseName(), e.After.GetPrimaryValues(e.Metadata.GetPrimaryKeys()))
}

func (e *NomiosEvent) GetPartitionID(size int) int {
	key := e.GetPartitionKey()

	h := fnv.New32a()
	h.Write([]byte(key))

	return int(h.Sum32()) % size
}
