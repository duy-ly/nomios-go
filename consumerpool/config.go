package consumerpool

import (
	"time"

	"github.com/spf13/viper"
)

type PoolConfig struct {
	PoolStreamSize int
	PoolSize       int
	BufferSize     int
	FlushTick      time.Duration
}

func loadConfig() PoolConfig {
	poolStreamSize := viper.GetInt("consumer.pool_stream_size")
	if poolStreamSize == 0 {
		poolStreamSize = 20000
	}
	poolSize := viper.GetInt("consumer.pool_size")
	if poolSize == 0 {
		poolSize = 5
	}
	bufferSize := viper.GetInt("consumer.buffer_size")
	if bufferSize == 0 {
		bufferSize = 1024
	}
	flushTick := viper.GetDuration("consumer.flush_tick")
	if flushTick == 0 {
		flushTick = 100 * time.Millisecond
	}

	return PoolConfig{
		PoolStreamSize: poolStreamSize,
		PoolSize:       poolSize,
		BufferSize:     bufferSize,
		FlushTick:      flushTick,
	}
}
