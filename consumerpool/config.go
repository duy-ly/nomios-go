package consumerpool

import "time"

type PoolConfig struct {
	Count      int
	BufferSize int
	FlushTick  time.Duration
}
