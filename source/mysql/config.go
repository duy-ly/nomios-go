package mysqlsource

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type MySQLSourceConfig struct {
	// mysql binlog config
	Addr             string
	User             string
	Pass             string
	Database         string
	TableIncludeList []string

	// syncer config
	ServerID uint32
}

func loadConfig() MySQLSourceConfig {
	return MySQLSourceConfig{
		Addr:             viper.GetString("source.mysql.addr"),
		User:             viper.GetString("source.mysql.user"),
		Pass:             viper.GetString("source.mysql.password"),
		Database:         viper.GetString("source.mysql.database"),
		TableIncludeList: viper.GetStringSlice("source.mysql.table_include_list"),
		ServerID:         viper.GetUint32("source.mysql.server_id"),
	}
}

func buildIncludeTableRegex(db string, tables []string) []string {
	regexes := make([]string, 0)

	if len(tables) == 0 {
		regexes = append(regexes, fmt.Sprintf("%s\\..*", db))
		return regexes
	}

	for _, t := range tables {
		regex := t
		if !strings.HasPrefix(regex, fmt.Sprintf("%s\\.", db)) {
			regex = fmt.Sprintf("%s\\.%s", db, regex)
		}
		if !strings.HasSuffix(regex, "$") {
			regex = fmt.Sprintf("%s$", regex)
		}

		regexes = append(regexes, regex)
	}

	return regexes
}
