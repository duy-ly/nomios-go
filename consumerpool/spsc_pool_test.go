package consumerpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type tsample struct {
}

func (t *tsample) hashcode(event any) int32 {
	return event.(int32)
}

var wg1 sync.WaitGroup

func TestNewSpscPool(t *testing.T) {
	var partitioner tsample
	pool := NewPoolWithFunc(2, 16, handleFunc1, &partitioner)
	pool.Start()

	for i := int32(0); i < 100; i++ {
		wg1.Add(1)
		pool.Submit(i)
	}

	wg1.Wait()

	pool.Stop()

	time.Sleep(time.Millisecond * 100)

}

func handleFunc1(events []any) {
	for _, event := range events {
		fmt.Println(event.(int32))
		wg1.Done()
	}
}
