package mysqlsource_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/duy-ly/nomios-go/model"
	mysqlsource "github.com/duy-ly/nomios-go/source/mysql"
	"github.com/duy-ly/nomios-go/testkit"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewMySQLSource(t *testing.T) {
	ctx := context.Background()

	viper.Set("source.mysql.user", "root")
	viper.Set("source.mysql.password", "12345678")
	viper.Set("source.mysql.database", "nomios_db")
	viper.Set("source.mysql.server_id", 100)

	type testCase struct {
		name           string
		sqlFile        string
		tableList      []string
		eventCount     int
		assertMaps     []map[string]interface{}
		assertFieldMap map[string]string
		assertNotMap   map[string]interface{}
	}

	suite := make([]testCase, 0)

	// suite = append(suite, testCase{
	// 	name:       "should_success_datetime",
	// 	sqlFile:    "testkit/resources/db/datetime_test.sql",
	// 	tableList:  []string{"datetime_test"},
	// 	eventCount: 1,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"timestamp_col":   time.Date(2022, 02, 02, 12, 30, 15, 0, time.Local),
	// 			"timestamp_col_2": time.Date(2022, 02, 02, 12, 30, 15, 300000000, time.Local),
	// 			"timestamp_col_3": time.Date(2022, 02, 02, 12, 30, 15, 222222000, time.Local),
	// 			"datetime_col":    time.Date(2022, 02, 02, 12, 30, 15, 0, time.UTC),
	// 			"datetime_col_2":  time.Date(2022, 02, 02, 12, 30, 15, 300000000, time.UTC),
	// 			"datetime_col_3":  time.Date(2022, 02, 02, 12, 30, 15, 222222000, time.UTC),
	// 			"year_col":        int64(2022),
	// 			"time_col":        time.Date(1970, 01, 01, 12, 30, 15, 0, time.UTC),
	// 			"date_col":        time.Date(2022, 02, 02, 0, 0, 0, 0, time.UTC),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_decimal",
	// 	sqlFile:    "testkit/resources/db/decimal_test.sql",
	// 	tableList:  []string{"nomios_db\\.DBZ730"},
	// 	eventCount: 1,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"A": 1.33,
	// 			"B": -2.111,
	// 			"C": 3.44400,
	// 			"D": 15.28,
	// 			"E": 0e-18,
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_default",
	// 	sqlFile:    "testkit/resources/db/default_value_test.sql",
	// 	tableList:  []string{"default_value"},
	// 	eventCount: 1,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"A": "a",
	// 			"B": 2.321,
	// 			"C": 999999,
	// 			"D": 0,
	// 			"E": int32(1),
	// 			"F": nil,
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_enum",
	// 	sqlFile:    "testkit/resources/db/enum_test.sql",
	// 	tableList:  []string{"test_stations"},
	// 	eventCount: 1,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"type": "station",
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_json",
	// 	sqlFile:    "testkit/resources/db/json_test.sql",
	// 	tableList:  []string{"dbz_126_jsontable"},
	// 	eventCount: 41,
	// 	assertMaps: make([]map[string]interface{}, 41),
	// 	assertFieldMap: map[string]string{
	// 		"json": "expectedBinlogStr",
	// 	},
	// 	assertNotMap: map[string]interface{}{
	// 		"id": nil,
	// 	},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_tinyint",
	// 	sqlFile:    "testkit/resources/db/tinyint_test.sql",
	// 	tableList:  []string{"DBZ1773"},
	// 	eventCount: 1,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"ti":  int64(100),
	// 			"ti1": int64(5),
	// 			"ti2": int64(50),
	// 			"ti3": int64(1),
	// 			"b":   int64(1),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_serial",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_1185_serial"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"id": int64(10),
	// 		},
	// 		{
	// 			"id": int64(11),
	// 		},
	// 		{
	// 			"id": int64(-1),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_serial_default",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_1185_serial_default_value"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"id": int64(10),
	// 		},
	// 		{
	// 			"id": int64(11),
	// 		},
	// 		{
	// 			"id": int64(1000),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_unsigned_int",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_228_int_unsigned"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"c1": int64(4294967295),
	// 			"c2": int64(4294967295),
	// 			"c3": int64(2147483647),
	// 			"c4": int64(4294967295),
	// 			"c5": int64(4294967295),
	// 			"c6": int64(2147483647),
	// 		},
	// 		{
	// 			"c1": int64(3294967295),
	// 			"c2": int64(3294967295),
	// 			"c3": int64(-1147483647),
	// 			"c4": int64(3294967295),
	// 			"c5": int64(3294967295),
	// 			"c6": int64(-1147483647),
	// 		},
	// 		{
	// 			"c1": int64(0),
	// 			"c2": int64(0),
	// 			"c3": int64(-2147483648),
	// 			"c4": int64(0),
	// 			"c5": int64(0),
	// 			"c6": int64(-2147483648),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_unsigned_smallint",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_228_smallint_unsigned"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"c1": int64(65535),
	// 			"c2": int64(65535),
	// 			"c3": int64(32767),
	// 		},
	// 		{
	// 			"c1": int64(45535),
	// 			"c2": int64(45535),
	// 			"c3": int64(-12767),
	// 		},
	// 		{
	// 			"c1": int64(0),
	// 			"c2": int64(0),
	// 			"c3": int64(-32768),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_unsigned_mediumint",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_228_mediumint_unsigned"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"c1": int64(16777215),
	// 			"c2": int64(16777215),
	// 			"c3": int64(8388607),
	// 		},
	// 		{
	// 			"c1": int64(10777215),
	// 			"c2": int64(10777215),
	// 			"c3": int64(-6388607),
	// 		},
	// 		{
	// 			"c1": int64(0),
	// 			"c2": int64(0),
	// 			"c3": int64(-8388608),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_unsigned_bigint",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_228_bigint_unsigned"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"c1": int64(-1),
	// 			"c2": int64(-1),
	// 			"c3": int64(9223372036854775807),
	// 		},
	// 		{
	// 			"c1": int64(-4000000000000000001),
	// 			"c2": int64(-4000000000000000001),
	// 			"c3": int64(-1223372036854775807),
	// 		},
	// 		{
	// 			"c1": int64(0),
	// 			"c2": int64(0),
	// 			"c3": int64(-9223372036854775808),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	// suite = append(suite, testCase{
	// 	name:       "should_success_unsigned_tinyint",
	// 	sqlFile:    "testkit/resources/db/unsigned_integer_test.sql",
	// 	tableList:  []string{"dbz_228_tinyint_unsigned"},
	// 	eventCount: 3,
	// 	assertMaps: []map[string]interface{}{
	// 		{
	// 			"c1": int64(255),
	// 			"c2": int64(255),
	// 			"c3": int64(127),
	// 		},
	// 		{
	// 			"c1": int64(155),
	// 			"c2": int64(155),
	// 			"c3": int64(-100),
	// 		},
	// 		{
	// 			"c1": int64(0),
	// 			"c2": int64(0),
	// 			"c3": int64(-128),
	// 		},
	// 	},
	// 	assertFieldMap: map[string]string{},
	// 	assertNotMap:   map[string]interface{}{},
	// })

	for _, tc := range suite {
		container := testkit.CreateMysqlContainer(filepath.Join(testkit.GetProjectRoot(), tc.sqlFile))
		defer container.Stop(ctx, new(time.Duration))

		port, err := container.MappedPort(ctx, "3306")
		if !assert.Equal(t, nil, err, "container started fail. tc %s", tc.name) {
			continue
		}

		viper.Set("source.mysql.addr", fmt.Sprintf("127.0.0.1:%d", port.Int()))
		viper.Set("source.mysql.table_include_list", tc.tableList)

		s, err := mysqlsource.NewMySQLSource()
		if !assert.Equal(t, nil, err, "source create fail. tc %s", tc.name) {
			continue
		}

		ch := make(chan []*model.NomiosEvent, 50)

		s.Start("", ch)

		time.Sleep(time.Second) // make sure source complete send all log message to channel

		s.Stop()

		time.Sleep(time.Second) // make sure source complete stop send new message to channel

		close(ch)

		assert.Equal(t, tc.eventCount, len(ch), "received event not equal. tc %s", tc.name)

		count := 0

		for events := range ch {
			e := events[0]

			count++

			assertMap := tc.assertMaps[count-1]
			for k, v := range assertMap {
				assert.Equal(t, v, e.After.Values[k], "data from source not equal. tc %s record %d field %s", tc.name, count, k)
			}

			for f1, f2 := range tc.assertFieldMap {
				s1, e1 := e.After.Values[f1].(string)
				s2, e2 := e.After.Values[f2].(string)
				if e1 && e2 {
					assert.JSONEq(t, s1, s2, "data field json from source not equal. tc %s record %d field 1 %s field 2 %s", tc.name, count, f1, f2)
				} else {
					assert.EqualValues(t, e.After.Values[f1], e.After.Values[f2], "data field from source not equal. tc %s record %d field 1 %s field 2 %s", tc.name, count, f1, f2)
				}
			}

			for k, v := range tc.assertNotMap {
				assert.NotEqual(t, v, e.After.Values[k], "data from source incorrect. tc %s record %d field %s", tc.name, count, k)
			}
		}
	}
}
