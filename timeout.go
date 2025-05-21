package onthewire

import (
	"fmt"
	"io"
	"time"
)

var ErrTimedOut = fmt.Errorf("operation timed out")

func conditionalAddTimeout(useTimeout bool, fn func([]byte) ([]byte, error), duration time.Duration) func([]byte) ([]byte, error) {
	if !useTimeout {
		return fn
	}

	return func(b []byte) ([]byte, error) {
		timeout := time.NewTimer(duration).C
		result := make(chan []byte, 1)
		err := make(chan error, 1)

		go func() {
			d, e := fn(b)
			if e != nil {
				err <- e
			} else {
				result <- d
			}
		}()

		select {
		case <-timeout:
			logger.Error("Failed to complete operation before timeout", "Error", ErrTimedOut, "Timeout", duration)
			return nil, ErrTimedOut
		case e := <-err:
			logger.Error("Failed to complete operation", "Error", e)
			return nil, e
		case data := <-result:
			return data, nil
		}
	}
}

func conditionalAddTimeoutReader(useTimeout bool, fn func(io.Reader) ([]byte, int, error), duration time.Duration) func(io.Reader) ([]byte, int, error) {
	if !useTimeout {
		return fn
	}

	return func(r io.Reader) ([]byte, int, error) {
		timeout := time.NewTimer(duration).C
		data := make(chan []byte, 1)
		count := make(chan int, 1)
		err := make(chan error, 1)

		go func() {
			d, c, e := fn(r)
			if e != nil {
				err <- e
			} else {
				data <- d
				count <- c
			}
		}()

		select {
		case <-timeout:
			logger.Error("Failed to complete read before timeout", "Error", ErrTimedOut, "Timeout", duration)
			return nil, 0, ErrTimedOut
		case e := <-err:
			logger.Error("Failed to complete read", "Error", e)
			return nil, 0, e
		case data := <-data:
			c := <-count
			return data, c, nil
		}
	}
}

func conditionalAddTimeoutWriter(useTimeout bool, fn func([]byte, io.Writer) (int, error), duration time.Duration) func([]byte, io.Writer) (int, error) {
	if !useTimeout {
		return fn
	}

	return func(b []byte, w io.Writer) (int, error) {
		timeout := time.NewTimer(duration).C
		count := make(chan int, 1)
		err := make(chan error, 1)

		go func() {
			c, e := fn(b, w)
			if e != nil {
				err <- e
			} else {
				count <- c
			}
		}()

		select {
		case <-timeout:
			logger.Error("Failed to complete write before timeout", "Error", ErrTimedOut, "Timeout", duration)
			return 0, ErrTimedOut
		case e := <-err:
			logger.Error("Failed to complete write", "Error", e)
			return 0, e
		case c := <-count:
			return c, nil
		}
	}
}
