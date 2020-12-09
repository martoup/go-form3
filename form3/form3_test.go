package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(server.URL, nil)
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func testQueryParam(t *testing.T, r *http.Request, key string, expectedValue string) {
	if r.URL.Query().Get(key) != expectedValue {
		t.Errorf("Request query param %v = %v, got %v", key, expectedValue, r.URL.Query().Get(key))
	}
}

func testBody(t *testing.T, r *http.Request, v io.Reader) {
	var expected, actual interface{}
	bodyDecoder := json.NewDecoder(r.Body)
	if err := bodyDecoder.Decode(&actual); err != nil {
		t.Errorf("Failed to compare objects. Decoding failed: %v", err)
	}

	expectedDecoder := json.NewDecoder(v)
	if err := expectedDecoder.Decode(&expected); err != nil {
		t.Errorf("Failed to compare objects. Decoding failed: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bodies not equal %+v, expected %+v", actual, expected)
	}
}

func TestNewClient_appendsV1(t *testing.T) {
	const baseURL = "http://localhost"
	c, err := NewClient(baseURL, nil)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	const want = baseURL + "/v1/"
	if got := c.baseURL.String(); got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClientFromEnvironment_readsEnv(t *testing.T) {
	os.Setenv("FORM3_BASE_URL", "http://localhost")
	defer os.Unsetenv("FORM3_BASE_URL")

	const baseURL = "http://localhost"
	c, err := NewClientFromEnvironment()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	const want = baseURL + "/v1/"
	if got := c.baseURL.String(); got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClientFromEnvironment_returnserr(t *testing.T) {
	_, err := NewClientFromEnvironment()
	if err == nil {
		t.Errorf("client should throw err")
	}
}

func TestCheckResponse_statusOk(t *testing.T) {
	r := http.Response{
		StatusCode: http.StatusOK,
	}
	err := CheckResponse(&r)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestCheckResponse_statusErr(t *testing.T) {
	parsedUrl, _ := url.Parse("http://localhost")
	r := http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"status:\"}:\"Internal Server Error\"}")),
		Request: &http.Request{
			Method: "GET",
			URL:    parsedUrl,
		},
	}
	err := CheckResponse(&r)
	if err == nil {
		t.Errorf("CheckResponse should throw an error.")
	}
}

func TestNewRequest(t *testing.T) {
	c, _ := NewClient("http://localhost", nil)
	wantUrl, _ := url.Parse("http://localhost/v1/test")
	request, err := c.NewRequest(http.MethodGet, "test", nil)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if got := request.Header.Get("Accept"); got != applicationJson {
		t.Errorf("Header is %v, want %v", got, applicationJson)
	}

	if got := request.Header.Get("Content-Type"); got != "" {
		t.Errorf("Header is not empty %v", got)
	}

	if got := request.URL.String(); got != wantUrl.String() {
		t.Errorf("Request URL is %v, want %v", got, wantUrl)
	}
}
