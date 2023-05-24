package consumerpool

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type test_event struct {
	x int32
	y int32
}

func (t *test_event) hashcode(event any) int32 {
	tmp := event.(*test_event)
	return tmp.x + tmp.y
}

func (t *test_event) compare(x, y any) int {

	tmpx := x.(*test_event)
	tmpy := y.(*test_event)

	return (int)(tmpx.x - tmpy.x)
}

var wg2 sync.WaitGroup

func TestNewSpscPool_Struct(t *testing.T) {
	var partitioner test_event
	pool := NewPoolWithFunc(2, 16, handleFunc2, &partitioner)
	pool.Start()

	for i := int32(0); i < 100; i++ {
		wg2.Add(1)
		pool.Submit(&test_event{i, rand.Int31()})

		if i%2 == 0 {
			last := pool.lastEvent()
			if last != nil {
				fmt.Println("last_event: ", pool.lastEvent().(*test_event))
			}
		}
	}

	wg2.Wait()

	pool.Stop()

	time.Sleep(time.Millisecond * 100)

	fmt.Println("last_event: ", pool.lastEvent().(*test_event))

}

func handleFunc2(events []any) {
	for _, event := range events {
		fmt.Println(event.(*test_event))
		wg2.Done()
	}

	time.Sleep(time.Millisecond * 5)
}

func TestCompare_test_truct(t *testing.T) {
	t1 := &test_event{}
	fmt.Println(t1.compare(&test_event{2, 2}, &test_event{3, 3}))

}
