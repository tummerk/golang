package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(
	baseURL string,
	httpClient *http.Client,
) APIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return APIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (a APIClient) Get(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	dest any,
	errDest any,
) (*http.Response, error) {
	return a.httpRequest(ctx, http.MethodGet, endpoint, headers, http.NoBody, dest, errDest)
}

func (a APIClient) Post(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	formData url.Values,
	dest any,
	errDest any,
) (*http.Response, error) {
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	body := bytes.NewBufferString(formData.Encode())

	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, body, dest, errDest)
}

func (a APIClient) PostJSON(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	requestJSON string,
	dest any,
	errDest any,
) (*http.Response, error) {
	b := []byte(requestJSON)

	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) MultiForm(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	payload io.Reader,
	dest any,
	errDest any,
) (*http.Response, error) {
	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, payload, dest, errDest)
}

func (a APIClient) httpRequest(
	ctx context.Context,
	httpMethod string,
	endpoint string,
	headers http.Header,
	payload io.Reader,
	dest any,
	errDest any,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, a.baseURL+endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	logRequest(req)

	if httpMethod == http.MethodPost && headers.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}

	defer resp.Body.Close()

	logResponse(resp)

	if err = parseResponse(resp, dest, errDest); err != nil {
		return nil, fmt.Errorf("parseResponse: %w", err)
	}

	return resp, nil
}

func parseResponse(r *http.Response, dest, errDest any) error {
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices && dest != nil {
		if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
			return fmt.Errorf("json.Decode(success destination): %w", err)
		}
	} else if errDest != nil {
		if err := json.NewDecoder(r.Body).Decode(errDest); err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("json.Decode(err destination): %w", err)
		}
	}

	return nil
}

func logRequest(r *http.Request) {
	log.Printf("Request:  %s %s", r.Method, r.URL)
}

func logResponse(r *http.Response) {
	rawResponse, err := httputil.DumpResponse(r, true)
	if err == nil {
		log.Println("Response:", string(rawResponse))
	}
}
