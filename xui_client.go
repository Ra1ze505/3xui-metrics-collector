package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type XUIClient struct {
	client  *http.Client
	config  *Config
	baseURL string
	cookies []*http.Cookie
}

type APIResponse struct {
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
	Obj     []Inbound `json:"obj"`
}

type Inbound struct {
	ID          int           `json:"id"`
	Up          int64         `json:"up"`
	Down        int64         `json:"down"`
	Total       int64         `json:"total"`
	Remark      string        `json:"remark"`
	Enable      bool          `json:"enable"`
	ExpiryTime  int64         `json:"expiryTime"`
	ClientStats []ClientStats `json:"clientStats"`
	Listen      string        `json:"listen"`
	Port        int           `json:"port"`
	Protocol    string        `json:"protocol"`
	Tag         string        `json:"tag"`
	Sniffing    string        `json:"sniffing"`
	Allocate    string        `json:"allocate"`
}

type ClientStats struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	ExpiryTime int64  `json:"expiryTime"`
	Total      int64  `json:"total"`
	Reset      int64  `json:"reset"`
}

func NewXUIClient(config *Config) *XUIClient {
	protocol := "http"
	if config.XUIUseTLS {
		protocol = "https"
	}

	return &XUIClient{
		client:  &http.Client{},
		config:  config,
		baseURL: fmt.Sprintf("%s://%s:%s%s", protocol, config.XUIHost, config.XUIPort, config.XUIBasePath),
	}
}

func (c *XUIClient) Login() error {
	data := url.Values{}
	data.Set("username", c.config.XUIUsername)
	data.Set("password", c.config.XUIPassword)

	loginURL := c.baseURL + "/login"
	log.Printf("Attempting to login at URL: %s", loginURL)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	c.cookies = resp.Cookies()
	log.Printf("Login successful, received %d cookies", len(c.cookies))
	return nil
}

func (c *XUIClient) GetInbounds() ([]Inbound, error) {
	inboundsURL := c.baseURL + "/panel/api/inbounds/list"
	log.Printf("Fetching inbounds from URL: %s", inboundsURL)

	req, err := http.NewRequest("GET", inboundsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create inbounds request: %w", err)
	}

	// Add cookies from login
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute inbounds request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get inbounds with status: %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode inbounds response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API returned error: %s", apiResp.Msg)
	}

	log.Printf("Successfully fetched %d inbounds", len(apiResp.Obj))
	return apiResp.Obj, nil
}
