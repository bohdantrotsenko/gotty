package server

import (
	"sync"
	"sync/atomic"
	"time"
)

type counter struct {
	duration  time.Duration
	zeroTimer *time.Timer
	wg        sync.WaitGroup

	// connections is the number of currently open connections.
	// The timer operations are serialized by mutex, the count is atomic.
	connections atomic.Int64
	mutex       sync.Mutex
}

func newCounter(duration time.Duration) *counter {
	zeroTimer := time.NewTimer(duration)

	// when duration is 0, drain the expire event here
	// so that user will never get the event.
	if duration == 0 {
		<-zeroTimer.C
	}

	return &counter{
		duration:  duration,
		zeroTimer: zeroTimer,
	}
}

func (counter *counter) add(n int) int {
	counter.mutex.Lock()
	if counter.duration > 0 {
		counter.zeroTimer.Stop()
	}
	counter.mutex.Unlock()

	counter.wg.Add(n)
	return int(counter.connections.Add(int64(n)))
}

func (counter *counter) done() int {
	n := int(counter.connections.Add(-1))
	counter.wg.Done()

	counter.mutex.Lock()
	if n == 0 && counter.duration > 0 {
		counter.zeroTimer.Reset(counter.duration)
	}
	counter.mutex.Unlock()

	return n
}

func (counter *counter) count() int {
	return int(counter.connections.Load())
}

func (counter *counter) wait() {
	counter.wg.Wait()
}

func (counter *counter) timer() *time.Timer {
	return counter.zeroTimer
}
