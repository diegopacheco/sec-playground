package mcp

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"
)

const (
	MethodInitialize       = "initialize"
	MethodInitialized      = "notifications/initialized"
	MethodToolsList        = "tools/list"
	MethodToolsCall        = "tools/call"
	MethodToolsListChanged = "notifications/tools/list_changed"
	MethodPing             = "ping"
)

type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name         string          `json:"name"`
	Title        string          `json:"title,omitempty"`
	Description  string          `json:"description,omitempty"`
	InputSchema  json.RawMessage `json:"inputSchema,omitempty"`
	OutputSchema json.RawMessage `json:"outputSchema,omitempty"`
	Annotations  json.RawMessage `json:"annotations,omitempty"`
}

type CallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

func (m Message) IsRequest() bool {
	return m.Method != "" && len(m.ID) > 0
}

func (m Message) IsNotification() bool {
	return m.Method != "" && len(m.ID) == 0
}

func (m Message) IsResponse() bool {
	return m.Method == "" && len(m.ID) > 0
}

func NewRequest(id any, method string, params any) (Message, error) {
	rawID, err := json.Marshal(id)
	if err != nil {
		return Message{}, err
	}
	msg := Message{JSONRPC: "2.0", ID: rawID, Method: method}
	if params != nil {
		msg.Params, err = json.Marshal(params)
	}
	return msg, err
}

func NewNotification(method string) Message {
	return Message{JSONRPC: "2.0", Method: method}
}

func NewResult(id json.RawMessage, result any) (Message, error) {
	raw, err := json.Marshal(result)
	return Message{JSONRPC: "2.0", ID: id, Result: raw}, err
}

func NewError(id json.RawMessage, code int, message string) Message {
	return Message{JSONRPC: "2.0", ID: id, Error: &Error{Code: code, Message: message}}
}

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReaderSize(r, 1<<16)}
}

func (r *Reader) ReadLine() ([]byte, error) {
	for {
		line, err := r.r.ReadBytes('\n')
		trimmed := trimLine(line)
		if len(trimmed) > 0 {
			return trimmed, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func trimLine(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r' || b[len(b)-1] == ' ') {
		b = b[:len(b)-1]
	}
	return b
}

type Writer struct {
	mu sync.Mutex
	w  io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteRaw(line []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	buf := make([]byte, 0, len(line)+1)
	buf = append(buf, line...)
	buf = append(buf, '\n')
	_, err := w.w.Write(buf)
	return err
}

func (w *Writer) Write(msg Message) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return w.WriteRaw(raw)
}
