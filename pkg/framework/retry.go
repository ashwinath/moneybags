package framework

import "time"

type retryOptions struct {
	maxRetries int
	delay      time.Duration
}

func (r *retryOptions) retry(fn func() error) error {
	var err error
	for i := 0; i < r.maxRetries; i++ {
		err = fn()
		if err == nil {
			break
		}

		time.Sleep(r.delay)
		r.delay *= 2
	}

	return err
}

func RetrySimple(fn func() error) error {
	ro := retryOptions{
		maxRetries: 3,
		delay:      30 * time.Second,
	}
	return ro.retry(fn)
}
