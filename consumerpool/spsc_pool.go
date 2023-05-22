package consumerpool

import (
	"sync"
	"sync/atomic"
)

// PoolWithFunc accepts the list of tasks from client,
// it limits the number of concurrency goroutine to Submit the input by an order.
type PoolWithFunc struct {
	// size of the pool: number of workers
	size int32

	// state iss used to notice the pool to close itself
	state int32

	// lock for protecting the status of queue.
	lock sync.Locker

	// poolFunc is the function for processing the list of tasks
	poolFunc func([]interface{})

	workers []*spsc_worker

	partitioner EventHelper

	wg *sync.WaitGroup
}

func NewPoolWithFunc(workerSize int32, bufferSize int32, pf func([]interface{}),
	partitioner EventHelper) *PoolWithFunc {
	if workerSize <= 0 {
		workerSize = 1
	}

	if pf == nil {
		return nil
	}

	p := &PoolWithFunc{
		size:        workerSize,
		state:       CREATED,
		lock:        new(sync.Mutex),
		poolFunc:    pf,
		partitioner: partitioner,
		wg:          new(sync.WaitGroup),
	}

	p.workers = make([]*spsc_worker, workerSize)

	for i := int32(0); i < workerSize; i++ {
		p.workers[i] = NewSpscWorker(bufferSize, p)
	}

	p.wg.Add(int(workerSize))

	return p
}

// Submit the event to run. This is a blocking function, if the correspond queue is full, it will wait here.
func (pool *PoolWithFunc) Submit(event any) error {

	if pool.status() != RUNNING {
		return ErrPoolStatusNotStarted
	}

	hashcode := pool.partitioner.hashcode(event)
	return pool.workers[hashcode%pool.size].Submit(event)
}

// lastEvent return the last committed event in all worker threads, it could be nil
// when the pool is running, but doest not submitted any message
func (pool *PoolWithFunc) lastEvent() any {
	var last_event any = nil
	for _, worker := range pool.workers {
		event := worker.lastEvent
		if event == nil {
			continue
		}

		if last_event == nil {
			last_event = event
			continue
		}

		// we have two different event, just compare by the interface, then return the oldest
		if pool.partitioner.compare(last_event, event) > 0 {
			last_event = event
		}

	}
	return last_event
}

func (pool *PoolWithFunc) Start() error {

	pool.lock.Lock()
	defer pool.lock.Unlock()

	//verify the state of the pool
	if pool.status() != CREATED {
		return ErrPoolStatusNotCreated
	}

	for _, worker := range pool.workers {
		worker.run()
	}

	atomic.StoreInt32(&pool.state, RUNNING)

	// lock then changed the state of a state, it's really should help
	return nil
}

func (pool *PoolWithFunc) Stop() {

	pool.lock.Lock()
	defer pool.lock.Unlock()

	atomic.StoreInt32(&pool.state, STOPPING)
	// from now on, the pool can not accept new tasks.

	for _, worker := range pool.workers {
		worker.Stop()
	}

	// wait for all worker is closed
	pool.wg.Wait()

	atomic.StoreInt32(&pool.state, STOPPED)
}

func (pool *PoolWithFunc) status() int32 {
	return atomic.LoadInt32(&pool.state)
}
