package resource_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func GetAllResourceProviders() ([]map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/resource-providers", Host), nil)
	if err != nil {
		return nil, err
	}

	r, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var response map[string][]map[string]string
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response["message"], nil
}

func WaitForCondition(condition func() bool, retries int, interval time.Duration) error {
	for i := 1; i <= retries; i++ {
		if condition() {
			return nil
		}

		time.Sleep(interval)
	}

	return errors.New("timeout")
}
