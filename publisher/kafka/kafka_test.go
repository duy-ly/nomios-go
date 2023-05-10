package kafkapublisher_test

import (
	"testing"

	"github.com/duy-ly/nomios-go/model"
	kafkapublisher "github.com/duy-ly/nomios-go/publisher/kafka"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/spf13/viper"
)

func Test_NewKafkaPublisher(t *testing.T) {
	viper.Set("publisher.host", "localhost")
	viper.Set("publisher.port", 9092)
	viper.Set("publisher.topic", "my-topic")

	p, err := kafkapublisher.NewKafkaPublisher()
	if err != nil {
		return
	}

	err = p.Publish([]*model.NomiosEvent{
		{
			Metadata: &model.MySQLMetadata{
				GTID:     "9ae7a6a4-ed69-11ed-9923-0242ac140004:4",
				ServerID: 112233,
				Ts:       1683530895,
				TableSchema: model.MySQLTableSchema{
					Name:     "test",
					Database: "catalog_db",
				},
				BinlogPosition: mysql.Position{
					Pos: 1163,
				},
			},
			Timestamp: 1683605994880,
			After: model.Row{
				Values: map[string]interface{}{
					"id":   1,
					"user": "YQ==",
				},
			},
		},
	})
	if err != nil {
		return
	}

	p.Close()
}
