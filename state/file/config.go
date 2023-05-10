package filestate

import "github.com/spf13/viper"

type FileStateConfig struct {
	Path string
}

func loadConfig() FileStateConfig {
	path := viper.GetString("state.file.path")
	if path == "" {
		path = "./checkpoint.nom"
	}

	return FileStateConfig{
		Path: path,
	}
}
