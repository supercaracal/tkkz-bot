package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	timeout = 300 * time.Second
	speaker = "human"
	ending  = "（ﾎﾞﾛﾝ"
)

var (
	regexpForMention  = regexp.MustCompile(`<@U[0-9A-Z]+>`)
	regexpForEmoticon = regexp.MustCompile(`:\S+:`)
	regexpForURL      = regexp.MustCompile(`<\S+>`)
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

	cli := buildClient()

	reply, err := fetchReply(cli, req)
	if err != nil {
		return err.Error()
	}

	return personalizeReply(trimReply(reply))
}

func buildRequest(apiURL, message string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request object for brain api: %w", err)
	}

	form := url.Values{}
	form.Add("speaker", speaker)
	form.Add("message", message)
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))

	return req, nil
}

func buildClient() *http.Client {
	client := http.Client{Timeout: timeout}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &client
}

func fetchReply(client *http.Client, req *http.Request) (string, error) {
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

func trimReply(text string) string {
	text = regexpForMention.ReplaceAllString(text, "")
	text = regexpForEmoticon.ReplaceAllString(text, "")
	text = regexpForURL.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "&lt; ", "")
	text = strings.ReplaceAll(text, "&gt; ", "")

	return strings.Trim(text, " 　\t\r\n")
}

func personalizeReply(text string) string {
	return fmt.Sprintf("%s%s", text, ending)
}
