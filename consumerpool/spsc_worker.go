package consumerpool

import "github.com/siddontang/go-log/log"

// spsc_worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks, then batching the tasks, then
// performing functions calls
type spsc_worker struct {

	// pool owns this worker
	pool *PoolWithFunc

	// args is a job should be done.
	args chan any

	lastEvent any
}

func NewSpscWorker(bufferSize int32, pool *PoolWithFunc) *spsc_worker {
	return &spsc_worker{
		pool: pool,
		args: make(chan any, bufferSize),
	}
}

func (w *spsc_worker) run() {
	go func() {
		defer func() {
			// when exit a running worker
			w.pool.wg.Done()
			log.Info("done one worker")
		}()

		for {
			select {
			case event := <-w.args:
				if event == nil {
					// this for closing the worker
					return
				}

				events := make([]any, 0)
				events = append(events, event)

				remainLen := len(w.args)
				for remainLen > 0 {
					remainLen--
					nextEvent := <-w.args
					if nextEvent == nil {
						// we still maybe have the nil (stopped signal) on the list
						return
					}
					events = append(events, nextEvent)
				}

				if w.pool.status() == RUNNING {
					// call Submit func for a batch of events (size at least is one)
					w.pool.poolFunc(events)

					w.lastEvent = events[len(events)-1]
				}
			}
		}

	}()
}

func (w *spsc_worker) Stop() {
	w.args <- nil
}

func (w *spsc_worker) Submit(arg interface{}) error {
	if w.pool.status() != RUNNING {
		return ErrPoolStatusNotStarted
	}
	w.args <- arg
	return nil
}
