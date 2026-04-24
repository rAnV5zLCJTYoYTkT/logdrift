package herald

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// WriterSink sends each message as a line to an io.Writer.
type WriterSink struct {
	w io.Writer
}

// NewWriterSink returns a WriterSink writing to w.
// If w is nil, os.Stderr is used.
func NewWriterSink(w io.Writer) *WriterSink {
	if w == nil {
		w = os.Stderr
	}
	return &WriterSink{w: w}
}

// Send writes msg followed by a newline.
func (ws *WriterSink) Send(msg string) error {
	_, err := fmt.Fprintln(ws.w, msg)
	return err
}

// WebhookSink POSTs each message as a plain-text body to a URL.
type WebhookSink struct {
	url    string
	client *http.Client
}

// NewWebhookSink returns a WebhookSink targeting url.
// Returns an error if url is empty.
func NewWebhookSink(url string) (*WebhookSink, error) {
	if url == "" {
		return nil, fmt.Errorf("herald: webhook url must not be empty")
	}
	return &WebhookSink{url: url, client: &http.Client{}}, nil
}

// Send POSTs msg to the configured URL.
func (ws *WebhookSink) Send(msg string) error {
	resp, err := ws.client.Post(ws.url, "text/plain", strings.NewReader(msg))
	if err != nil {
		return fmt.Errorf("herald: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("herald: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
