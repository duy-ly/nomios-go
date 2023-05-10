package mysqlsource

import "github.com/spf13/viper"

type MySQLSourceConfig struct {
	// mysql binlog config
	Addr string
	User string
	Pass string

	// syncer config
	ServerID uint32
}

func loadConfig() MySQLSourceConfig {
	return MySQLSourceConfig{
		Addr:     viper.GetString("source.mysql.addr"),
		User:     viper.GetString("source.mysql.user"),
		Pass:     viper.GetString("source.mysql.password"),
		ServerID: viper.GetUint32("source.mysql.server_id"),
	}
}
