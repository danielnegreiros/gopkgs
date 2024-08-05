package toolpkgs

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {

	client := NewHTTPClient(10, 5, 5)

	httpReq := HTTPRequest{
		Client:        client,
		Timeout:       10,
		EndPoint:      "http://httpbin.org/get",
		Method:        "GET",
		Body:          []byte(`{"key": "value"}`),
		Data:          &struct{}{},
		AcceptedCodes: []int{200},
	}

	contentBytes, statusCode, _ := httpReq.Execute()
	if statusCode != httpReq.AcceptedCodes[0] {
		t.Errorf("\nStatus Code. Expected: %d, Found %d", httpReq.AcceptedCodes[0], statusCode)
	}
	log.Println(contentBytes)
}

func TestPost(t *testing.T) {

	client := NewHTTPClient(10, 5, 5)
	httpReq := HTTPRequest{
		Client:        client,
		Timeout:       10,
		EndPoint:      "http://httpbin.org/post",
		Method:        "POST",
		Body:          []byte(`{"key": "value"}`),
		Data:          &struct{}{},
		AcceptedCodes: []int{200},
	}

	contentBytes, statusCode, _ := httpReq.Execute()
	if statusCode != httpReq.AcceptedCodes[0] {
		t.Errorf("\nStatus Code. Expected: %d, Found %d", httpReq.AcceptedCodes[0], statusCode)
	}
	log.Println(contentBytes)
}
