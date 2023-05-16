package consumerpool

import "errors"

var (
	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrPoolStatusNotCreated will be return when you want to start a pool, but it's not on status CREATED
	ErrPoolStatusNotCreated = errors.New("pool status not CREATED, you can not take next action")

	// ErrPoolStatusNotStarted will be return when you want to start a pool, but it's not on status CREATED
	ErrPoolStatusNotStarted = errors.New("pool status not RUNNING, you can not take next action")
)
