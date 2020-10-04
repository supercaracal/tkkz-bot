package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetDefaultReply is
func GetDefaultReply(apiURL, message string) string {
	if apiURL == "" || message == "" {
		return ""
	}

	req, err := buildRequest(apiURL, message)
	if err != nil {
		return err.Error()
	}

	reply, err := fetchReply(req)
	if err != nil {
		return err.Error()
	}

	return reply
}

func buildRequest(apiURL, message string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request object for brain api: %w", err)
	}

	q := req.URL.Query()
	q.Add("speaker", "human")
	q.Add("message", message)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func fetchReply(req *http.Request) (string, error) {
	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to request to brain api: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || 300 <= res.StatusCode {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("Failed to read response body for brain api: %w", err)
		}

		return "", fmt.Errorf("brain api replied %d: %s", res.StatusCode, string(body))
	}

	reply, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body for brain api: %w", err)
	}

	return string(reply), nil
}
