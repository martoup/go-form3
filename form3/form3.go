package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	// applicationJson holds value that can be used in HTTP headers
	// in requests manipulating JSON data
	applicationJson = "application/json; charset=utf-8"
	baseURLKey      = "FORM3_BASE_URL"
)

// Client manages the communication with the Form3 API.
type Client struct {
	// HTTP client used to communicate with the Form3 API.
	httpClient *http.Client
	// Base URL for API requests.
	baseURL *url.URL
	// Accounts holds a reference to an AccountService
	// which handles the communication with the account related methods of the Form3 API.
	Accounts *AccountsService
}

// service is a type that holds a reference to a Client and allows unified way of managing services.
type service struct {
	client *Client
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred.
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context should not be nil")
	}
	req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("WARNING: could not close body %v", err)
		}
	}()

	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	requestURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, requestURL.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", applicationJson)
	}
	req.Header.Set("Accept", applicationJson)
	return req, nil
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.
// Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.ErrorMessage = string(data)
		}
	}

	return errorResponse
}

// ErrorResponse represents an error caused by an API request
type ErrorResponse struct {
	Response     *http.Response // HTTP response that caused this error
	ErrorMessage string         `json:"error_message"`
}

// Error returns a string representation of an API request error
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.ErrorMessage)
}

// NewClient returns a new Form3 API client, using the given
// http.Client to perform all requests. If a nil httpClient is
// provided, a new http.Client will be used.
// BaseURL to the Form3 API should be provided in the format http(s)://host:port
//
// Users who wish to pass their own http.Client should use this method.
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}
	if !strings.HasSuffix(parsedURL.Path, "/v1/") {
		parsedURL.Path += "v1/"
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    parsedURL,
	}
	c.Accounts = &AccountsService{client: c}
	return c, nil
}

// NewClientFromEnvironment returns a new Form3 API client.
// The baseURL to the API will be read from the environment and default http client will be used.
// Users who wish to pass their own http.Client should use NewClient() method.
func NewClientFromEnvironment() (*Client, error) {
	baseURL, exists := os.LookupEnv(baseURLKey)
	if !exists {
		msg := fmt.Sprintf("Please set the baseURL %s environment variable", baseURLKey)
		return nil, errors.New(msg)
	}
	return NewClient(baseURL, nil)
}
