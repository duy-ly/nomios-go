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

	sm := state.NewStateManager(s, mockCollectFn)
	sm.Start()

	time.Sleep(3 * time.Second)

	sm.Stop()

	time.Sleep(time.Millisecond) // make sure state cron complete stop

	assert.Equal(t, "30", sm.GetLastCheckpoint(), "check state manager cron")

	sm.Checkpoint()

	assert.Equal(t, "31", sm.GetLastCheckpoint(), "check state manager manual checkpoint")
}
