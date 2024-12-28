package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type Client struct {
	Url               string
	Client            *http.Client
	ReconcileRetries  int
	ReconcileInterval time.Duration
}

func NewClient(url string, timeout time.Duration, reconcileReties int, reconcileInterval time.Duration) Client {
	return Client{
		Url:               url,
		Client:            &http.Client{Timeout: timeout},
		ReconcileRetries:  reconcileReties,
		ReconcileInterval: reconcileInterval,
	}
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

func Complement(condition func() bool) func() bool {
	return func() bool {
		return !condition()
	}
}

func (c Client) Reconcile(entity string, name string, url string) func() bool {
	return func() bool {
		allEntities, err := c.FetchAll(entity)
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

func (c Client) FetchAll(entity string) ([]map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%ss", c.Url, entity), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.Client.Do(req)
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

func (c Client) Post(entity string, name string, url string) error {
	postBody, err := json.Marshal(map[string]string{
		"name": name,
		"url":  url,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%ss", c.Url, entity),
		bytes.NewBuffer(postBody),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = c.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) Delete(entity string, name string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/%ss/%s", c.Url, entity, name),
		nil,
	)

	_, err = c.Client.Do(req)

	return err
}

func (c Client) ReconcilePipeline(pipeline map[string]any) func() bool {
	return func() bool {
		expected, err := c.FetchPipeline(pipeline["group"].(string), pipeline["name"].(string))

		return err == nil && reflect.DeepEqual(expected, pipeline)
	}
}

func (c Client) FetchPipeline(group string, name string) (map[string]any, error) {
	params := url.Values{
		"group": {group},
		"name":  {name},
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/pipelines?%s", c.Url, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	r, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var response map[string][]map[string]any
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	pipelines := response["message"]
	if len(pipelines) == 0 {
		return nil, fmt.Errorf("no such pipeline")
	}

	return pipelines[0], nil
}

func (c Client) PostPipeline(attrs map[string]any) error {
	postBody, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		c.Url+"%s/pipelines",
		bytes.NewBuffer(postBody),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = c.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) DeletePipeline(group string, name string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/pipelines/groups/%s/names/%s", c.Url, group, name),
		nil,
	)

	_, err = c.Client.Do(req)

	return err
}

func (c Client) ReconcilePipelineDeletion(group string, name string) func() bool {
	return func() bool {
		_, err := c.FetchPipeline(group, name)

		return err != nil && err.Error() == "no such pipeline"
	}
}
