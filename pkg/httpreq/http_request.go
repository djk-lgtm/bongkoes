package httpreq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient struct {
	endpoint string
	client   *http.Client
	username string
	password string
}

type Opts struct {
	Endpoint string
	Username string
	Password string
}

func NewHTTPClient(o *Opts) *HTTPClient {
	client := http.DefaultClient
	return &HTTPClient{
		client:   client,
		endpoint: o.Endpoint,
		username: o.Username,
		password: o.Password,
	}
}

func (c *HTTPClient) ExecuteBasicAuth(ctx context.Context, method string, path string, headers map[string]string, jsonBytes []byte) ([]byte, error) {
	req, err := http.NewRequest(method, c.endpoint+path, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	req = req.WithContext(ctx)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		fmt.Println(fmt.Sprintf("[HTTP Client] response:%+v", string(resBody)))
		return nil, errors.New("non success http call")
	}

	return resBody, nil
}
