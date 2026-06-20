package xui

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string, timeout time.Duration, insecure bool) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL: baseURL,
		token:   token,
		http: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		},
	}
}

type envelope struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     json.RawMessage `json:"obj"`
}

func (c *Client) do(ctx context.Context, method, path string, out any) error {
	url := c.baseURL + path
	var body io.Reader
	if method == http.MethodPost {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var env envelope
	if err := json.Unmarshal(respBody, &env); err != nil {
		return fmt.Errorf("decode envelope: %w", err)
	}
	if !env.Success {
		if env.Msg != "" {
			return fmt.Errorf("API error: %s", env.Msg)
		}
		return fmt.Errorf("API error: success=false")
	}
	if out == nil {
		return nil
	}
	if len(env.Obj) == 0 || string(env.Obj) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Obj, out); err != nil {
		return fmt.Errorf("decode obj: %w", err)
	}
	return nil
}

func (c *Client) GetServerStatus(ctx context.Context) (*ServerStatus, error) {
	var status ServerStatus
	if err := c.do(ctx, http.MethodGet, "/panel/api/server/status", &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (c *Client) GetNodes(ctx context.Context) ([]Node, error) {
	var nodes []Node
	if err := c.do(ctx, http.MethodGet, "/panel/api/nodes/list", &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *Client) GetInbounds(ctx context.Context) ([]Inbound, error) {
	var inbounds []Inbound
	if err := c.do(ctx, http.MethodGet, "/panel/api/inbounds/list", &inbounds); err != nil {
		return nil, err
	}
	return inbounds, nil
}

func (c *Client) GetClients(ctx context.Context) ([]ClientWithAttachments, error) {
	var clients []ClientWithAttachments
	if err := c.do(ctx, http.MethodGet, "/panel/api/clients/list", &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *Client) GetOnlines(ctx context.Context) ([]string, error) {
	var onlines []string
	if err := c.do(ctx, http.MethodPost, "/panel/api/clients/onlines", &onlines); err != nil {
		return nil, err
	}
	return onlines, nil
}

func (c *Client) GetLastOnline(ctx context.Context) (map[string]int64, error) {
	var lastOnline map[string]int64
	if err := c.do(ctx, http.MethodPost, "/panel/api/clients/lastOnline", &lastOnline); err != nil {
		return nil, err
	}
	if lastOnline == nil {
		lastOnline = map[string]int64{}
	}
	return lastOnline, nil
}

func (c *Client) GetOutboundsTraffic(ctx context.Context) ([]OutboundTraffic, error) {
	var outbounds []OutboundTraffic
	if err := c.do(ctx, http.MethodGet, "/panel/api/xray/getOutboundsTraffic", &outbounds); err != nil {
		return nil, err
	}
	return outbounds, nil
}

// ParseEnvelope decodes a raw API response body (for tests).
func ParseEnvelope(body []byte, out any) error {
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}
	if !env.Success {
		if env.Msg != "" {
			return fmt.Errorf("API error: %s", env.Msg)
		}
		return fmt.Errorf("API error: success=false")
	}
	if out == nil {
		return nil
	}
	if len(env.Obj) == 0 || string(env.Obj) == "null" {
		return nil
	}
	return json.Unmarshal(env.Obj, out)
}
