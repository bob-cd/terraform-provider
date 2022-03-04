package common

import (
	"errors"
	"net/http"
	"time"
)

var (
	Host   = "http://localhost:7777"
	Client = &http.Client{Timeout: 10 * time.Second}
)

func WaitForCondition(condition func() bool, retries int, interval time.Duration) error {
	for i := 1; i <= retries; i++ {
		if condition() {
			return nil
		}

		time.Sleep(interval)
	}

	return errors.New("timeout")
}
