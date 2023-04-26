package consumerpool

import "github.com/duy-ly/nomios-go/event"

type ConsumerPool struct {
	stopSig chan bool

	consumers []*Consumer
}

// NewConsumerPool -- create a pool of consumer
func NewConsumerPool(count int) *ConsumerPool {
	p := new(ConsumerPool)
	p.stopSig = make(chan bool)

	for i := 0; i <= count; i++ {
		c := NewConsumer(i)

		c.Start()

		p.consumers = append(p.consumers, c)
	}

	return p
}

// Start -- start consumer pool, get partition using hash function to pick consumer then handle NomiosEvent
func (p *ConsumerPool) Start(stream chan event.NomiosEvent) {
	go func() {
		for {
			select {
			case <-p.stopSig:
				break
			case e := <-stream:
				// TODO: do partition
				partitionIdx := 1

				p.consumers[partitionIdx].Send(e)
			}
		}
	}()
}
