package state

import (
	"time"

	"github.com/spf13/viper"
)

type StateManagerConfig struct {
	Cron time.Duration
}

func loadConfig() StateManagerConfig {
	return StateManagerConfig{
		Cron: viper.GetDuration("state.checkpoint_cron"),
	}
}
