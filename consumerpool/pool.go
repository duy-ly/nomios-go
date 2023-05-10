package consumerpool

import (
	"sync"
	"time"

	"github.com/duy-ly/nomios-go/model"
	"github.com/duy-ly/nomios-go/util"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type ConsumerPool struct {
	cfg     PoolConfig
	stopSig chan bool
	flushed chan bool

	consumers []*Consumer
}

// NewConsumerPool -- create a pool of consumer
func NewConsumerPool(cfg PoolConfig) *ConsumerPool {
	if cfg.Count <= 0 {
		cfg.Count = 1
	}
	if cfg.BufferSize <= 0 {
		cfg.BufferSize = 100
	}
	if cfg.FlushTick <= 0 {
		cfg.FlushTick = 100 * time.Millisecond
	}

	p := new(ConsumerPool)
	p.stopSig = make(chan bool, 1)
	p.flushed = make(chan bool, 1)

	for i := 0; i < cfg.Count; i++ {
		c := NewConsumer(i, cfg.BufferSize, cfg.FlushTick)

		c.Start()

		p.consumers = append(p.consumers, c)
	}

	return p
}

// Start -- start consumer pool, get partition using hash function to pick consumer then handle NomiosEvent
func (p *ConsumerPool) Start(stream chan []*model.NomiosEvent) {
	go func() {
		for {
			select {
			case <-p.stopSig:
				// ensure get all event from stream
				for e := range stream {
					p.partitionEvent(e)
				}

				p.flushed <- true
				return
			case e := <-stream:
				p.partitionEvent(e)
			}
		}
	}()
}

func (p *ConsumerPool) GetLastEventPos() string {
	var minPos *mysql.Position

	for _, c := range p.consumers {
		e := c.lastProcessedEvent.Load()
		if e == nil {
			continue
		}

		ePos := e.Metadata.GetPos()

		if minPos == nil || minPos.Compare(ePos) == 1 {
			minPos = &ePos
		}
	}

	if minPos == nil {
		return ""
	}

	return util.BuildEventPos(minPos)
}

func (p *ConsumerPool) partitionEvent(e []*model.NomiosEvent) {
	// TODO: do partition
	partitionIdx := 1

	p.consumers[partitionIdx].Send(e)
}

func (p *ConsumerPool) Stop() {
	p.stopSig <- true
	<-p.flushed

	var wg sync.WaitGroup
	for _, c := range p.consumers {
		wg.Add(1)
		go func(c *Consumer) {
			defer wg.Done()
			c.Stop()
		}(c)
	}
	wg.Wait()
}
