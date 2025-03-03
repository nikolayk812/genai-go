package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func NewLoggingRoundTripper(transport http.RoundTripper) *LoggingRoundTripper {
	return &LoggingRoundTripper{Transport: transport}
}

func (c *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := logRequest(req); err != nil {
		return nil, fmt.Errorf("logRequest: %w", err)
	}

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("RoundTrip: %w", err)
	}

	if err := logResponse(resp); err != nil {
		return nil, fmt.Errorf("logResponse: %w", err)
	}

	return resp, nil
}

func logRequest(req *http.Request) error {
	log.Printf("Request: %s %s", req.Method, req.URL)

	if req.Body == nil {
		return nil
	}

	var err error
	req.Body, err = logBody(req.Body, "Request Body", true)
	return err
}

func logResponse(resp *http.Response) error {
	return nil

	//log.Printf("Response Status: %s", resp.Status)
	//
	//if resp.Body == nil {
	//	return nil
	//}
	//
	//var err error
	//resp.Body, err = logBody(resp.Body, "Response Body", false)
	//return err
}

func logBody(body io.ReadCloser, logPrefix string, pretty bool) (io.ReadCloser, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	printBytes := bodyBytes
	if pretty {
		// Final response for SSE is not valid JSON
		// {"model":"llama3.2","created_at":"2025-03-03T12:46:54.866443Z","message":{"role":"assistant","content":"How"},"done":false}
		// {"model":"llama3.2","created_at":"2025-03-03T12:46:54.904178Z","message":{"role":"assistant","content":" can"},"done":false}
		// {"model":"llama3.2","created_at":"2025-03-03T12:46:55.177874Z","message":{"role":"assistant","content":""},"done_reason":"stop","done":true,"total_duration":943410791,"load_duration":48098083,"prompt_eval_count":26,"prompt_eval_duration":582000000,"eval_count":8,"eval_duration":311000000}

		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err != nil {
			return nil, fmt.Errorf("json.Indent: %w", err)
		}
		printBytes = prettyJSON.Bytes()
	}

	log.Printf("%s: %s", logPrefix, string(printBytes))

	// Reset the body so it can be read again
	return io.NopCloser(bytes.NewBuffer(bodyBytes)), nil
}
