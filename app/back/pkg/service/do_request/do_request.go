package do_request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
	token    string
}

func New(host string, basePath string) Client {
	return Client{
		host:     host,
		basePath: basePath,
		client:   http.Client{},
	}
}

func (c *Client) TokenSet(token string) {
	c.token = token
}

func (c *Client) DoRequestPost(method string, jsonBody []byte, httpMethod string, query url.Values) ([]byte, int, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(httpMethod, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 404, fmt.Errorf("%s cannot do request 1: %w", method, err)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}

	res, err1 := c.client.Do(req)
	if err1 != nil {
		return nil, 404, fmt.Errorf("%s cannot do request 2: %w", method, err1)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode > 300 {
		body, _ := io.ReadAll(res.Body)
		// todo логирование fmt.Errorf("%s Non-OK HTTP status: %d, BODY: %s", method, res.StatusCode, string(body))
		return body, res.StatusCode, fmt.Errorf("%s Non-OK HTTP status: %d", method, res.StatusCode)
	}

	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return nil, 404, fmt.Errorf("%s cannot do request 3: %w", method, err2)
	}

	return body, 200, nil
}

func (c *Client) DoRequestSendFile(method, contentType string, body bytes.Buffer) ([]byte, int, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), &body)
	if err != nil {
		return nil, 404, fmt.Errorf("%s cannot do request 1: %w", method, err)
	}
	req.Header.Add("Content-Type", contentType)

	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}

	res, err1 := c.client.Do(req)
	if err1 != nil {
		return nil, 404, fmt.Errorf("%s cannot do request 2: %w", method, err1)
	}
	defer func() { _ = res.Body.Close() }()

	bodyRes, _ := io.ReadAll(res.Body)
	if res.StatusCode > 300 {
		fmt.Println("Non-OK HTTP status: ", res.StatusCode)

		return nil, res.StatusCode, fmt.Errorf("%s Non-OK HTTP status: %d, BODY: %s", method, res.StatusCode, string(bodyRes))
	}

	return bodyRes, 200, nil
}

func (c *Client) DoRequest(method string, query url.Values) ([]byte, int, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, 404, fmt.Errorf("cannot do request 1: %w", err)
	}

	req.URL.RawQuery = query.Encode()
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	res, err1 := c.client.Do(req)
	if err1 != nil {
		sCode := 0
		if res != nil {
			sCode = res.StatusCode
		}
		return nil, sCode, fmt.Errorf("cannot do request 2: %w", err1)
	}
	defer func() { _ = res.Body.Close() }()
	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return nil, 404, fmt.Errorf("cannot do request 3: %w", err2)
	}

	return body, res.StatusCode, nil
}
