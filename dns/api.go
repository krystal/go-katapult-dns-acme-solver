package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) baseHost() string {
	if c.Host == "" {
		return "https://api.katapult.io"
	}

	return c.Host
}

func (c *Client) apiRequest(method string, baseURL string, params map[string]string, requestBody string) ([]byte, error) {
	// Parse the base URL
	u, err := url.Parse(fmt.Sprintf("%s/%s", c.baseHost(), baseURL))
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	// Convert the params into URL values
	q := url.Values{}
	for key, value := range params {
		q.Add(key, value)
	}

	// Add the encoded parameters to the URL
	u.RawQuery = q.Encode()

	// Create a new request
	var req *http.Request
	if requestBody == "" {
		req, err = http.NewRequest(method, u.String(), nil)
	} else {
		req, err = http.NewRequest(method, u.String(), strings.NewReader(requestBody))
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	// Add the authorization header to the request
	req.Header.Add("Authorization", "Bearer "+c.APIToken)

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	// If the result is not JSON, return an error.
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("body was not in JSON format, not appropriate")
	}

	// If the status is OK, we'll return that body now for future
	// string parsing.
	if resp.StatusCode == http.StatusOK {
		return body, nil
	}

	// Otherwise, we'll need to parse this as an error
	errorResponse := &ErrorResponse{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal error response: %w", err)
	}

	// We now just need to turn our error response in to an actual error
	// which can be returned and look at.
	return body, fmt.Errorf("error: %s (%s)", errorResponse.Error.Code, errorResponse.Error.Description)
}
