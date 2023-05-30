package breaker

import (
	"errors"
	"fmt"
	"time"
)

// Settings configures CircuitBreaker:
//
// Name is the name of the CircuitBreaker.
//
// MaxRequests is the maximum number of requests allowed to pass through
// when the CircuitBreaker is half-open.
// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
//
// Interval is the cyclic period of the closed state
// for the CircuitBreaker to clear the internal Counts.
// If Interval is less than or equal to 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
//
// Timeout is the period of the open state,
// after which the state of the CircuitBreaker becomes half-open.
// If Timeout is less than or equal to 0, the timeout value of the CircuitBreaker is set to 60 seconds.
//
// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
// If ReadyToTrip is nil, default ReadyToTrip is used.
// Default ReadyToTrip returns true when the number of consecutive failures is more than 5.
//
// OnStateChange is called whenever the state of the CircuitBreaker changes.
//
// IsSuccessful is called with the error returned from a request.
// If IsSuccessful returns true, the error is counted as a success.
// Otherwise the error is counted as a failure.
// If IsSuccessful is nil, default IsSuccessful is used, which returns false for all non-nil errors.
type setting struct {
	name          string
	maxRequests   uint32
	interval      time.Duration
	timeout       time.Duration
	readyToTrip   func(counts Counts) bool
	onStateChange func(name string, from State, to State)
	isSuccessful  func(err error) bool
}

func newSetting() *setting {
	return &setting{
		name:        "unspecified-breaker",
		maxRequests: 1,
		interval:    time.Duration(0) * time.Second,
		timeout:     time.Duration(60) * time.Second,
		readyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		isSuccessful: func(err error) bool {
			return err == nil
		},
	}
}

func (s *setting) apply(opts []Option) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return fmt.Errorf("could not apply setting: %w", err)
		}
	}

	return nil
}

type Option func(*setting) error

func WithIsSuccessful(input func(err error) bool) Option {
	return func(s *setting) error {
		if input == nil {
			return errors.New("can not set isSuccessful as nil")
		}

		s.isSuccessful = input

		return nil
	}
}

func WithonStateChange(input func(name string, from State, to State)) Option {
	return func(s *setting) error {
		if input == nil {
			return errors.New("can not set onStateChange as nil")
		}
		s.onStateChange = input

		return nil
	}
}

func WithReadyToTrip(input func(counts Counts) bool) Option {
	return func(s *setting) error {
		if input == nil {
			return errors.New("can not set readyToTrip as nil")
		}
		s.readyToTrip = input

		return nil
	}
}

func WithTimeout(input time.Duration) Option {
	return func(s *setting) error {
		if input <= 0 {
			return nil
		}

		s.timeout = input

		return nil
	}
}

func WithName(name string) Option {
	return func(s *setting) error {
		if name == "" {
			return nil
		}

		s.name = name

		return nil
	}
}

func WithMaxRequests(input uint32) Option {
	return func(s *setting) error {
		if input <= 1 {
			return errors.New("can not set max request less than 1")
		}

		s.maxRequests = input

		return nil
	}
}

func WithInterval(input time.Duration) Option {
	return func(s *setting) error {
		if input <= 0 {
			return nil
		}

		s.interval = input

		return nil
	}
}
