package utils

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreaker implements a circuit breaker pattern to handle failure cases.
type CircuitBreaker struct {
	maxFailures   int
	timeout       time.Duration
	failureCount  int
	lastFailure   time.Time
	state         string
	mutex         sync.Mutex
}

// NewCircuitBreaker initializes a circuit breaker.
func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       "closed",
	}
}

// Execute runs the function and monitors failures.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == "open" {
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = "half-open"
		} else {
			return errors.New("circuit breaker is open")
		}
	}

	err := fn()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()

		if cb.failureCount >= cb.maxFailures {
			cb.state = "open"
		}
		return err
	}

	if cb.state == "half-open" {
		cb.state = "closed"
	}
	cb.failureCount = 0
	return nil
}

// SimpleCircuitBreaker is a simpler implementation of a circuit breaker.
type SimpleCircuitBreaker struct {
	failureThreshold int
	failureCount     int
	lastFailureTime  time.Time
	retryTimeout     time.Duration
}

// NewSimpleCircuitBreaker initializes a simple circuit breaker.
func NewSimpleCircuitBreaker(failureThreshold int, retryTimeout time.Duration) *SimpleCircuitBreaker {
	return &SimpleCircuitBreaker{
		failureThreshold: failureThreshold,
		retryTimeout:     retryTimeout,
	}
}

// Call runs the function and monitors failures.
func (cb *SimpleCircuitBreaker) Call(fn func() error) error {
	if cb.failureCount >= cb.failureThreshold {
		if time.Since(cb.lastFailureTime) < cb.retryTimeout {
			return errors.New("circuit breaker is open")
		}
		cb.failureCount = 0
	}

	err := fn()
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		return err
	}

	cb.failureCount = 0
	return nil
}
