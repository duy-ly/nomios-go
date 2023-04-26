package consumerpool

import (
	"sync"
	"time"

	"github.com/duy-ly/nomios-go/event"
)

type Consumer struct {
	mu          sync.Mutex
	bufferSize  int
	bufferQueue []event.NomiosEvent
	flushSig    chan bool
	publisher   *Publisher

	// state
	lastProcessedEvent *event.NomiosEvent
}

func NewConsumer(partition int) *Consumer {
	c := new(Consumer)
	c.bufferSize = 100
	c.bufferQueue = make([]event.NomiosEvent, c.bufferSize)
	c.flushSig = make(chan bool, 1)

	c.publisher = NewPublisher()

	return c
}

func (c *Consumer) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				c.handle()
			case <-c.flushSig:
				c.handle()
			}
		}
	}()
}

func (c *Consumer) GetLastProcessedEvent() *event.NomiosEvent {
	return c.lastProcessedEvent
}

func (c *Consumer) Send(e event.NomiosEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bufferQueue = append(c.bufferQueue, e)
	if len(c.bufferQueue) == c.bufferSize {
		c.flushSig <- true
	}
}

func (c *Consumer) handle() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tmp := make([]event.NomiosEvent, c.bufferSize)
	tmp, c.bufferQueue = c.bufferQueue, tmp

	err := c.publisher.Publish(tmp)
	if err != nil {
		return err
	}
	c.lastProcessedEvent = &tmp[len(tmp)-1]

	return nil
}
