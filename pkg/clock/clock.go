package clock

import (
	"errors"
	"sync"
	"time"

	clockapi "github.com/benbjohnson/clock"
)

var (
	globalClock = clockapi.New()
	mu          sync.Mutex
)

// Set the global clock to a clockapi.Mock with the given time.Time
func Set(t time.Time) {
	mu.Lock()
	defer mu.Unlock()
	mock, ok := globalClock.(*clockapi.Mock)
	if !ok {
		mock = clockapi.NewMock()
	}
	mock.Set(t)
	globalClock = mock
}

// Add moves the mocked global clock forward the given duration. It will error
// if the global clock is not mocked.
func Add(d time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	mock, ok := globalClock.(*clockapi.Mock)
	if !ok {
		return errors.New("time not mocked")
	}
	mock.Add(d)
	return nil
}

// Reset sets the global clock to a pure time implementation
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	globalClock = clockapi.New()
}

// Clock is a non-package level wrapper around time that supports stubbing.
// It will use its localized stubs (allowing for parallelized unit tests
// where package level stubbing would cause issues). It falls back to any
// package level time stubs for non-parallel, cross-package integration
// testing scenarios.
//
// If nothing is stubbed, it defaults to default time behavior in the time
// package.
type Clock struct {
	mock *clockapi.Mock
	sync.RWMutex
}

// Set sets the Clock to a clock.Mock at the given time.Time
func (c *Clock) Set(t time.Time) {
	c.Lock()
	defer c.Unlock()
	if c.mock == nil {
		c.mock = clockapi.NewMock()
	}
	c.mock.Set(t)
}

// Add moves clock forward time.Duration if it is mocked. It will error
// if the clock is not mocked.
func (c *Clock) Add(d time.Duration) error {
	c.Lock()
	defer c.Unlock()
	if c.mock == nil {
		return errors.New("clock not mocked")
	}
	c.mock.Add(d)
	return nil
}

// Reset removes local clock.Mock
func (c *Clock) Reset() {
	c.Lock()
	defer c.Unlock()
	c.mock = nil
}

func (c *Clock) After(d time.Duration) <-chan time.Time {
	if c.mock == nil {
		return globalClock.After(d)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.After(d)
}

func (c *Clock) AfterFunc(d time.Duration, f func()) *clockapi.Timer {
	if c.mock == nil {
		return globalClock.AfterFunc(d, f)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.AfterFunc(d, f)
}

func (c *Clock) Now() time.Time {
	if c.mock == nil {
		return globalClock.Now()
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.Now()
}

func (c *Clock) Since(t time.Time) time.Duration {
	if c.mock == nil {
		return globalClock.Since(t)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.Since(t)
}

func (c *Clock) Sleep(d time.Duration) {
	if c.mock == nil {
		globalClock.Sleep(d)
		return
	}
	c.RLock()
	defer c.RUnlock()
	c.mock.Sleep(d)
}

func (c *Clock) Tick(d time.Duration) <-chan time.Time {
	if c.mock == nil {
		return globalClock.Tick(d)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.Tick(d)
}

func (c *Clock) Ticker(d time.Duration) *clockapi.Ticker {
	if c.mock == nil {
		return globalClock.Ticker(d)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.Ticker(d)
}

func (c *Clock) Timer(d time.Duration) *clockapi.Timer {
	if c.mock == nil {
		return globalClock.Timer(d)
	}
	c.RLock()
	defer c.RUnlock()
	return c.mock.Timer(d)
}
