package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// loggingTrasport is a transport that logs the request and response.
// It's useful for debugging the HTTP client to be used with the LLM.
type loggingTrasport struct{}

// RoundTrip implements the http.RoundTripper interface, logging the request and response bodies.
func (c *loggingTrasport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("RoundTrip: %s", req.URL)

	// read the body
	body, err := io.ReadAll(req.Body)
	req.Body.Close() // close the original body
	if err != nil {
		return nil, err
	}
	// log the response
	log.Printf("Request: %s", string(body))

	// create a new ReadCloser with the same content
	req.Body = io.NopCloser(bytes.NewReader(body))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	log.Printf("Status: %s", resp.Status)

	// read the body
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close() // close the original body
	if err != nil {
		return resp, err
	}

	// log the body
	log.Printf("Response: %s", string(body))

	// create a new ReadCloser with the same content
	resp.Body = io.NopCloser(bytes.NewReader(body))

	return resp, err
}
