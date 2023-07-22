package utils

import (
	"fmt"
	"time"
)

// FN represents the logic to be retried.
type FN func() (bool, error)

// Do executes a function with a generic retry mechanism.
func Do(tries int, timeBetweenTries time.Duration, fn FN) error {
	var err error
	var keepTrying bool

	for attempt := 1; ; attempt++ {
		keepTrying, err = fn()
		if err == nil {
			return nil
		}

		if !keepTrying {
			return err
		}

		if attempt >= tries {
			break
		}

		time.Sleep(timeBetweenTries * time.Duration(attempt))
	}

	return fmt.Errorf("giving up after %d attempts, last error: %w", tries, err)
}
