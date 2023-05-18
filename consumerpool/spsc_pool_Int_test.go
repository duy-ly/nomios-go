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

func (t *tsample) compare(x, y any) int {
	return int(x.(int32) - y.(int32))
}

var wg1 sync.WaitGroup

func TestNewSpscPool_Int(t *testing.T) {
	var partitioner tsample
	pool := NewPoolWithFunc(2, 16, handleFunc1, &partitioner)
	pool.Start()

	for i := int32(0); i < 100; i++ {
		wg1.Add(1)
		pool.Submit(i)

		if i%2 == 0 {
			last := pool.lastEvent()
			if last != nil {
				fmt.Println("last_event: ", pool.lastEvent().(int32))
			}
		}
	}

	wg1.Wait()

	pool.Stop()

	time.Sleep(time.Millisecond * 100)

	fmt.Println("last_event: ", pool.lastEvent().(int32))

}

func handleFunc1(events []any) {
	for _, event := range events {
		fmt.Println(event.(int32))
		wg1.Done()
	}

	time.Sleep(time.Millisecond * 5)
}

func TestCompare(t *testing.T) {
	t1 := &tsample{}
	fmt.Println(t1.compare(int32(2), int32(3)))

}
