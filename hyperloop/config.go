package hyperloop

import (
	"github.com/spf13/viper"
)

type HyperloopConfig struct {
	EventStreamSize int
	StateKind       string
	SourceKind      string
	PublisherKind   string
}

func loadConfig() HyperloopConfig {
	streamSize := viper.GetInt("hyperloop.stream_size")
	if streamSize == 0 {
		streamSize = 20000
	}
	stateKind := viper.GetString("hyperloop.state")
	if stateKind == "" {
		stateKind = "file"
	}
	sourceKind := viper.GetString("hyperloop.source")
	if sourceKind == "" {
		sourceKind = "mysql"
	}
	publisherKind := viper.GetString("hyperloop.publisher")
	if publisherKind == "" {
		publisherKind = "dummy"
	}

	return HyperloopConfig{
		EventStreamSize: streamSize,
		StateKind:       stateKind,
		SourceKind:      sourceKind,
		PublisherKind:   publisherKind,
	}
}
