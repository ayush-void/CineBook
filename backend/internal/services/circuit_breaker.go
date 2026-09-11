package services

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failures         int
	failureThreshold int
	cooldown         time.Duration
	nextTry          time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: threshold,
		cooldown:         cooldown,
	}
}

// Execute runs action through the breaker, tripping to Open after
// failureThreshold consecutive failures and resetting after cooldown.
func (cb *CircuitBreaker) Execute(action func() error) error {
	cb.mu.Lock()
	if cb.state == Open {
		if time.Now().After(cb.nextTry) {
			cb.state = HalfOpen
		} else {
			cb.mu.Unlock()
			return errors.New("circuit open: payment gateway unavailable")
		}
	}
	cb.mu.Unlock()

	err := action()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		if cb.failures >= cb.failureThreshold {
			cb.state = Open
			cb.nextTry = time.Now().Add(cb.cooldown)
		}
		return err
	}

	cb.failures = 0
	cb.state = Closed
	return nil
}
