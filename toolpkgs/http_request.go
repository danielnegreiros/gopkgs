package toolpkgs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// NewHTTPClient creates a new HTTP client with specified timeout and connection settings.
func NewHTTPClient(timeout int, maxIdleConn int, maxIdleConnPerHost int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        maxIdleConn,
			MaxIdleConnsPerHost: maxIdleConnPerHost,
			DisableKeepAlives:   false,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		},
	}
}

// ComposeCredentials encodes user credentials into a URL-encoded form.
func ComposeCredentials(user string, pass string) []byte {
	body := url.Values{}
	body.Set("username", user)
	body.Add("password", pass)
	return []byte(body.Encode())
}

// HTTPRequest represents an HTTP request with various configuration options.
type HTTPRequest struct {
	Client        *http.Client
	Timeout       int
	EndPoint      string
	Method        string
	Body          []byte
	AcceptedCodes []int
	Data          interface{}
	Cookie        *http.Cookie
	Header        map[string]string
}

// Execute performs the HTTP request and returns the response data and status code.
func (r *HTTPRequest) Execute() (interface{}, int, error) {
	req, err := http.NewRequest(r.Method, r.EndPoint, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, 0, err
	}

	if r.Cookie != nil {
		req.AddCookie(r.Cookie)
	}

	for k, v := range r.Header {
		req.Header.Add(k, v)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	log.Printf("%s %s %d", req.Method, req.URL.Path, resp.StatusCode)
	if !isAcceptedCode(resp.StatusCode, r.AcceptedCodes) {
		log.Println(string(contentBytes))
		return nil, resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(contentBytes, r.Data)
	if err != nil {
		return nil, 0, err
	}

	return r.Data, resp.StatusCode, nil
}

// isAcceptedCode checks if the HTTP status code is in the list of accepted codes.
func isAcceptedCode(code int, acceptedCodes []int) bool {
	for _, allow := range acceptedCodes {
		if allow == code {
			return true
		}
	}
	return false
}
