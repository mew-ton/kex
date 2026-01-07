package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RemoteProvider struct {
	BaseURL string // URL to directory (or wherever kex.json is relative to)
	KexURL  string // Full URL to kex.json
	Token   string // Optional Bearer Token
}

func NewRemoteProvider(rootURL, token string) *RemoteProvider {
	baseURL := rootURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	kexURL := baseURL + "kex.json"

	return &RemoteProvider{
		BaseURL: baseURL,
		KexURL:  kexURL,
		Token:   token,
	}
}

func (r *RemoteProvider) Load() (*IndexSchema, []error) {
	req, err := http.NewRequest("GET", r.KexURL, nil)
	if err != nil {
		return nil, []error{err}
	}

	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, []error{err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, []error{fmt.Errorf("failed to fetch kex.json: status %d", resp.StatusCode)}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []error{err}
	}

	schema := &IndexSchema{}
	if err := json.Unmarshal(data, schema); err != nil {
		return nil, []error{err}
	}

	return schema, nil
}

func (r *RemoteProvider) FetchContent(path string) (string, error) {
	url := path
	if !strings.HasPrefix(path, "http") {
		url = r.BaseURL + strings.TrimLeft(path, "/")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sContent := string(body)
	parts := strings.SplitN(sContent, "\n---\n", 2)
	if len(parts) >= 2 {
		return parts[1], nil
	}
	return sContent, nil
}
