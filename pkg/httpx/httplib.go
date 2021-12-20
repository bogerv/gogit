package httpx

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type request struct {
	headers map[string]string
}

func New() *request {
	return &request{
		headers: make(map[string]string),
	}
}

func (r *request) AddHeader(key, value string) *request {
	r.headers[key] = value
	return r
}

func Get(url string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if resp == nil {
		return "", errors.New("req is nil")
	}

	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func (r *request) Post(url string, data interface{}) ([]byte, error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Add("content-type", "application/json;charset=UTF-8")
	if r.headers != nil {
		for key, value := range r.headers {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}
