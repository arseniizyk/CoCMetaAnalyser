package cocdev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const (
	baseURL = "https://developer.clashofclans.com/api"
)

type Client struct {
	client   *http.Client
	email    string
	password string
}

type Key struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	CidrRanges  []string `json:"cidrRanges"`
}

type listResponse struct {
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Keys []Key `json:"keys"`
}

type createResponse struct {
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Key Key `json:"key"`
}

func New(email, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Jar:     jar,
	}

	return &Client{
		client:   client,
		email:    email,
		password: password,
	}, nil
}

func (c *Client) Login() error {
	payload := map[string]string{
		"email":    c.email,
		"password": c.password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/login", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

func (c *Client) GetPublicIP() (string, error) {
	endpoints := []string{
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
		"https://icanhazip.com",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	for _, ep := range endpoints {
		resp, err := client.Get(ep)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			ipBytes, err := io.ReadAll(resp.Body)
			if err == nil {
				ip := strings.TrimSpace(string(ipBytes))
				if ip != "" {
					return ip, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to determine public IP address")
}

func (c *Client) ListKeys() ([]Key, error) {
	req, err := http.NewRequest("POST", baseURL+"/apikey/list", bytes.NewBufferString("{}"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("listing keys failed: %w", err)
	}
	defer resp.Body.Close()

	var result listResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding key list: %w", err)
	}

	return result.Keys, nil
}

func (c *Client) CleanupOldKeys(prefix string) error {
	keys, err := c.ListKeys()
	if err != nil {
		return err
	}

	for _, k := range keys {
		if strings.HasPrefix(k.Name, prefix) {
			_ = c.RevokeKey(k.ID)
		}
	}
	return nil
}

func (c *Client) CreateKey(name, ip string) (*Key, error) {
	payload := map[string]any{
		"name":        name,
		"description": "Auto-generated temporary key for CoCMetaAnalyser",
		"cidrRanges":  []string{ip},
		"scopes":      []string{"clash"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+"/apikey/create", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("creating key failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("creating key failed with status %d: %s", resp.StatusCode, string(respBytes))
	}

	var result createResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding create key response: %w", err)
	}

	return &result.Key, nil
}

func (c *Client) RevokeKey(keyID string) error {
	payload := map[string]string{
		"id": keyID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/apikey/revoke", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("revoking key failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
