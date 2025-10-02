package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	httpMethodGet    = "get"
	httpMethodPost   = "post"
	httpMethodPut    = "put"
	httpMethodDelete = "delete"
)

func httpGet(u string, v interface{}) error {
	resp, code, err := httpCall("", u, nil)

	if !successful(code) {
		return fmt.Errorf("status code: %d: error: %v", code, err)
	}

	b, _ := io.ReadAll(strings.NewReader(*resp))
	err = json.Unmarshal(b, &v)

	return err
}
func httpPost(u string, payload, response interface{}) error {
	return httpWithPayload(httpMethodPost, u, payload, response)
}

func httpPut(u string, payload, response interface{}) error {
	return httpWithPayload(httpMethodPut, u, payload, response)
}
func httpDelete(u string, response interface{}) error {
	return httpWithPayload(httpMethodDelete, u, nil, response)
}

func httpWithPayload(method, u string, payload, response interface{}) (err error) {
	resp, _, err := httpCall(method, u, payload)
	if err != nil {
		return
	}

	b, _ := io.ReadAll(strings.NewReader(*resp))

	err = json.Unmarshal(b, &response)

	return
}
func httpCall(method, u string, payload interface{}) (resp *string, code int, err error) {
	p := map[string]any{}
	if payload != nil {
		v, e := json.Marshal(payload)
		buf := bytes.NewBuffer(v)
		if e != nil {
			err = e
			return
		}
		p["body"] = string(buf.Bytes())
	}

	if method == "" {
		method = httpMethodGet
	}
	p["method"] = method
	p["headers"] = map[string]any{
		"Content-Type": "application/json",
	}
	p["redirect"] = "error"

	resp, code, err = fetch(u, &p)

	return
}

func successful(code int) bool {
	return 100 <= code && code <= 399
}
