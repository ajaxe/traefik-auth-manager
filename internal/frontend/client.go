package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpGet(u string, v interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error code: %v", res.StatusCode)
	}

	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &v)

	return err
}
func httpPost(u string, payload, response interface{}) error {
	return httpWithPayload(http.MethodPost, u, payload, response)
}

func httpPut(u string, payload, response interface{}) error {
	return httpWithPayload(http.MethodPut, u, payload, response)
}
func httpDelete(u string, response interface{}) error {
	return httpWithPayload(http.MethodDelete, u, nil, response)
}

func httpWithPayload(method, u string, payload, response interface{}) error {
	buf := bytes.NewBuffer([]byte{})

	if payload != nil {
		v, err := json.Marshal(payload)
		buf = bytes.NewBuffer(v)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u, buf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &response)

	return err
}
