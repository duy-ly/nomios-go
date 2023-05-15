package state_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/duy-ly/nomios-go/state"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockState struct {
	counter int32
}

func (m *mockState) Reset() {
	m.counter = 0
}

func (m *mockState) SaveLastID(string) {
	atomic.AddInt32(&m.counter, 1)
}

func (m *mockState) GetLastID() string {
	return fmt.Sprintf("%d", m.counter)
}

func mockCollectFn() string {
	return ""
}

func Test_NewStateManager(t *testing.T) {
	viper.Set("state.checkpoint_cron", 100*time.Millisecond)

	s := &mockState{}

	before := func() {
		s.Reset()
	}

	type testCase struct {
		name           string
		waitBeforeStop time.Duration
		collectFn      func() string
		expectCron     string
		expectManual   string
	}

	suite := make([]testCase, 0)

	suite = append(suite, testCase{
		name: "should_success_update_state",

		waitBeforeStop: 3 * time.Second,
		collectFn:      mockCollectFn,
		expectCron:     "30",
		expectManual:   "31",
	})

	suite = append(suite, testCase{
		name:         "should_success_no_collect_no_update",
		expectCron:   "0",
		expectManual: "0",
	})

	for _, tc := range suite {
		before()

		sm := state.NewStateManager(s, tc.collectFn)
		sm.Start()

		time.Sleep(tc.waitBeforeStop)

		sm.Stop()

		time.Sleep(time.Millisecond) // make sure state cron complete stop

		assert.Equal(t, tc.expectCron, sm.GetLastCheckpoint(), "check state manager cron. tc %s", tc.name)

		sm.Checkpoint()

		assert.Equal(t, tc.expectManual, sm.GetLastCheckpoint(), "check state manager manual checkpoint. tc %s", tc.name)
	}
}
