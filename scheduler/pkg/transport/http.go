package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"scheduler/internal/logger"
	"strings"
	"sync"
	"time"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type HttpRequest struct {
	HttpClient *http.Client
	Url        string
	Port       int
	Method     string
	Headers    map[string]string
	Body       []byte
	Query      map[string]string
	Params     map[string]string
}

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Timeout:   time.Second * 300,
		Transport: transport,
	}
}

func MakeHTTPRequest(ctx context.Context, req HttpRequest) (*http.Response, error) {
	httpClient := req.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if httpClient.Timeout == 0 {
		httpClient.Timeout = 30 * time.Second
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	buf.ReadFrom(bytes.NewReader(req.Body))

	for {
		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Url, buf)
		if err != nil {
			return nil, err
		}

		for key, value := range req.Headers {
			httpReq.Header.Add(key, value)
		}

		if req.Query != nil {
			query := httpReq.URL.Query()
			for key, value := range req.Query {
				query.Add(key, value)
			}
			httpReq.URL.RawQuery = query.Encode()
		}

		if req.Params != nil {
			path := httpReq.URL.Path
			for key, value := range req.Params {
				path = strings.Replace(path, ":"+key, value, -1)
			}
			httpReq.URL.Path = path
		}

		res, err := httpClient.Do(httpReq)
		if err != nil {
			return nil, err
		}

		logger.Log.Debug(fmt.Sprintf("HTTP response status from %s is %s", req.Url, res.Status))
		return res, nil
	}
}

func RequestAndParseJSONBody(ctx context.Context, req HttpRequest) error {
	resp, err := MakeHTTPRequest(ctx, req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logger.Log.Debug(fmt.Sprintf("response body from %s is %s", req.Url, body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("unexpected status: " + resp.Status)
	}

	return err
}
