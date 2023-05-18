package consumerpool

const (
	// CREATED represents that the pool is created, not started yet, so that it can not accept new tasks now
	CREATED = iota

	// RUNNING represents that the pool is started, can accept new tasks
	RUNNING

	// STOPPING represents that the pool is beginning stopping, can not accept new tasks.
	STOPPING

	// STOPPED represents that the pool is stopped, can not accept new tasks.
	STOPPED
)

type EventHelper interface {

	// hashcode hashes the event into an int32 value
	hashcode(event any) int32

	// compare to event x, and y. If x is greater than y, return a positive value.
	// If x and y are equal, return zero.
	// other will return negative value
	compare(x, y any) int
}
