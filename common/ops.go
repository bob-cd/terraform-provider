package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func Reconcile(entity string, name string, url string) func() bool {
	return func() bool {
		allEntities, err := FetchAll(entity)
		if err != nil {
			return false
		}

		for _, entity := range allEntities {
			if entity["name"] == name && entity["url"] == url {
				return true
			}
		}

		return false
	}
}

func FetchAll(entity string) ([]map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%ss", Host, entity), nil)
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

func Post(entity string, name string, url string) error {
	postBody, _ := json.Marshal(map[string]string{
		"url": url,
	})

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%ss/%s", Host, entity, name),
		bytes.NewBuffer(postBody),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func Delete(entity string, name string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/%ss/%s", Host, entity, name),
		nil,
	)

	_, err = Client.Do(req)

	return err
}
