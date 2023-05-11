package consumerpool

import (
	"sync"

	"github.com/duy-ly/nomios-go/model"
)

type ConsumerPool struct {
	cfg     PoolConfig
	stopSig chan bool
	flushed chan bool
	stream  chan []*model.NomiosEvent

	consumers []*Consumer
}

// NewConsumerPool -- create a pool of consumer
func NewConsumerPool() (*ConsumerPool, error) {
	cfg := loadConfig()

	p := new(ConsumerPool)
	p.stopSig = make(chan bool, 1)
	p.flushed = make(chan bool, 1)
	p.stream = make(chan []*model.NomiosEvent, cfg.PoolStreamSize)

	for i := 0; i < cfg.PoolSize; i++ {
		c, err := NewConsumer(i, cfg.BufferSize, cfg.FlushTick)
		if err != nil {
			return nil, err
		}

		c.Start()

		p.consumers = append(p.consumers, c)
	}

	return p, nil
}

// Start -- start consumer pool, get partition using hash function to pick consumer then handle NomiosEvent
func (p *ConsumerPool) Start() {
	go func() {
		for {
			select {
			case <-p.stopSig:
				close(p.stream)

				// ensure get all event from stream
				for e := range p.stream {
					p.partitionEvent(e)
				}

				p.flushed <- true
				return
			case e := <-p.stream:
				p.partitionEvent(e)
			}
		}
	}()
}

func (p *ConsumerPool) GetStream() chan []*model.NomiosEvent {
	return p.stream
}

func (p *ConsumerPool) GetLastGTID() string {
	var lastEvent *model.NomiosEvent

	for _, c := range p.consumers {
		e := c.lastProcessedEvent.Load()
		if e == nil {
			continue
		}

		if lastEvent == nil {
			lastEvent = e
			continue
		}

		ePos := e.Metadata.GetPos()
		lastPos := lastEvent.Metadata.GetPos()

		if lastPos.Compare(ePos) == 1 {
			lastEvent = e
		}
	}

	if lastEvent == nil {
		return ""
	}

	return lastEvent.Metadata.GetID()
}

func (p *ConsumerPool) partitionEvent(events []*model.NomiosEvent) {
	mapPartitionEvents := make(map[int][]*model.NomiosEvent)

	for _, e := range events {
		partitionID := e.GetPartitionID(len(p.consumers))

		if _, exist := mapPartitionEvents[partitionID]; !exist {
			mapPartitionEvents[partitionID] = make([]*model.NomiosEvent, 0)
		}

		mapPartitionEvents[partitionID] = append(mapPartitionEvents[partitionID], e)
	}

	for partitionID, e := range mapPartitionEvents {
		p.consumers[partitionID].Send(e)
	}
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
