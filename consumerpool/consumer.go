package consumerpool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/publisher"
)

type Consumer struct {
	mu          sync.Mutex
	partition   int
	bufferSize  int
	flushTick   time.Duration
	bufferQueue [][]*model.NomiosEvent
	flushSig    chan bool
	stopSig     chan bool
	publisher   publisher.Publisher

	// state
	lastProcessedEvent atomic.Pointer[model.NomiosEvent]
}

func NewConsumer(partition int, bufferSize int, flushTick time.Duration) (*Consumer, error) {
	p, err := publisher.NewPublisher()
	if err != nil {
		return nil, err
	}

	c := new(Consumer)
	c.partition = partition
	c.bufferSize = bufferSize
	c.flushTick = flushTick
	c.bufferQueue = make([][]*model.NomiosEvent, 0)
	c.flushSig = make(chan bool, 1)
	c.stopSig = make(chan bool, 1)
	c.publisher = p

	c.lastProcessedEvent = atomic.Pointer[model.NomiosEvent]{}

	return c, err
}

func (c *Consumer) Start() {
	go func() {
		ticker := time.NewTicker(c.flushTick)

		for {
			select {
			case <-c.stopSig:
				ticker.Stop()
				return
			case <-ticker.C:
				c.handle()
			case <-c.flushSig:
				c.handle()
			}
		}
	}()
}

func (c *Consumer) GetLastProcessedEvent() *model.NomiosEvent {
	return c.lastProcessedEvent.Load()
}

func (c *Consumer) Send(e []*model.NomiosEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bufferQueue = append(c.bufferQueue, e)
	if len(c.bufferQueue) >= c.bufferSize {
		c.flushSig <- true
	}
}

func (c *Consumer) handle() error {
	if len(c.bufferQueue) == 0 {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	tmp := make([][]*model.NomiosEvent, 0)
	tmp, c.bufferQueue = c.bufferQueue, tmp

	for _, e := range tmp {
		if len(e) == 0 {
			continue
		}

		err := c.publisher.Publish(e)
		if err != nil {
			return err
		}

		c.lastProcessedEvent.Store(e[len(e)-1])
	}

	return nil
}

func (c *Consumer) Stop() {
	c.stopSig <- true
	c.handle()
	c.publisher.Close()
}
