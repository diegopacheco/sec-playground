package sim

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"mcp-rug-pull-proxy/internal/mcp"
)

type Client struct {
	w             *mcp.Writer
	incoming      chan mcp.Message
	mu            sync.Mutex
	notifications []string
	seq           int
	Timeout       time.Duration
}

func NewClient(fromServer io.Reader, toServer io.Writer) *Client {
	c := &Client{w: mcp.NewWriter(toServer), incoming: make(chan mcp.Message, 64), Timeout: 10 * time.Second}
	go func() {
		reader := mcp.NewReader(fromServer)
		for {
			line, err := reader.ReadLine()
			if err != nil {
				close(c.incoming)
				return
			}
			var msg mcp.Message
			if json.Unmarshal(line, &msg) == nil {
				c.incoming <- msg
			}
		}
	}()
	return c
}

func (c *Client) Request(method string, params any) (mcp.Message, error) {
	c.mu.Lock()
	c.seq++
	id := c.seq
	c.mu.Unlock()
	req, err := mcp.NewRequest(id, method, params)
	if err != nil {
		return mcp.Message{}, err
	}
	if err := c.w.Write(req); err != nil {
		return mcp.Message{}, err
	}
	want := string(req.ID)
	deadline := time.After(c.Timeout)
	for {
		select {
		case msg, ok := <-c.incoming:
			if !ok {
				return mcp.Message{}, errors.New("server closed the connection")
			}
			if msg.IsNotification() {
				c.record(msg.Method)
				continue
			}
			if msg.IsResponse() && string(msg.ID) == want {
				return msg, nil
			}
		case <-deadline:
			return mcp.Message{}, fmt.Errorf("timeout waiting for %s", method)
		}
	}
}

func (c *Client) Notify(method string) error {
	return c.w.Write(mcp.NewNotification(method))
}

func (c *Client) WaitNotification(method string, timeout time.Duration) bool {
	if c.seen(method) {
		return true
	}
	deadline := time.After(timeout)
	for {
		select {
		case msg, ok := <-c.incoming:
			if !ok {
				return false
			}
			if msg.IsNotification() {
				c.record(msg.Method)
				if msg.Method == method {
					return true
				}
			}
		case <-deadline:
			return false
		}
	}
}

func (c *Client) Initialize() (string, error) {
	msg, err := c.Request(mcp.MethodInitialize, map[string]any{
		"protocolVersion": "2025-06-18",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "client-sim", "version": "1.0.0"},
	})
	if err != nil {
		return "", err
	}
	if msg.Error != nil {
		return "", errors.New(msg.Error.Message)
	}
	var result struct {
		ServerInfo struct {
			Name string `json:"name"`
		} `json:"serverInfo"`
	}
	_ = json.Unmarshal(msg.Result, &result)
	return result.ServerInfo.Name, c.Notify(mcp.MethodInitialized)
}

func (c *Client) ListTools() ([]mcp.Tool, error) {
	msg, err := c.Request(mcp.MethodToolsList, map[string]any{})
	if err != nil {
		return nil, err
	}
	if msg.Error != nil {
		return nil, errors.New(msg.Error.Message)
	}
	var result struct {
		Tools []mcp.Tool `json:"tools"`
	}
	if err := json.Unmarshal(msg.Result, &result); err != nil {
		return nil, err
	}
	return result.Tools, nil
}

func (c *Client) CallTool(name string, args map[string]any) (string, *mcp.Error, error) {
	msg, err := c.Request(mcp.MethodToolsCall, map[string]any{"name": name, "arguments": args})
	if err != nil {
		return "", nil, err
	}
	if msg.Error != nil {
		return "", msg.Error, nil
	}
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	_ = json.Unmarshal(msg.Result, &result)
	if len(result.Content) == 0 {
		return "", nil, nil
	}
	return result.Content[0].Text, nil, nil
}

func (c *Client) record(method string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notifications = append(c.notifications, method)
}

func (c *Client) seen(method string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, m := range c.notifications {
		if m == method {
			return true
		}
	}
	return false
}
