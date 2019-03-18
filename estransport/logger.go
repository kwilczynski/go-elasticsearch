package estransport

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Logger represents the default logger.
//
type Logger struct {
	output io.Writer
	format string
}

// NewLogger returns new logger, when w is not nil, otherwise it returns nil.
//
func NewLogger(w io.Writer, f string) *Logger {
	if w == nil {
		return nil
	}
	return &Logger{output: w, format: f}
}

func (l *Logger) logRoundTrip(req *http.Request, res *http.Response, dur time.Duration) {
	fmt.Fprintf(l.output, "%s %s %s [status:%d request:%s]\n",
		time.Now().Format(time.RFC3339),
		req.Method,
		req.URL.String(),
		res.StatusCode,
		dur.Truncate(time.Millisecond),
	)
	l.logRequestBody(req)
}

func (l *Logger) logRequestBody(req *http.Request) {
	if req.Body != nil && req.Body != http.NoBody {
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			for _, line := range strings.Split(string(body), "\n") {
				if line != "" {
					fmt.Fprintf(l.output, "> %s\n", line)
				}
			}
		}
	}
}

func (l *Logger) logResponseBody(res *http.Response) {
	if res.Body != nil && res.Body != http.NoBody {
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			defer func() { res.Body = ioutil.NopCloser(bytes.NewReader(body)) }()
			defer func() { res.Body.Close() }()
			for _, line := range strings.Split(string(body), "\n") {
				if line != "" {
					fmt.Fprintf(l.output, "< %s\n", line)
				}
			}
		}
	}
}

func (l *Logger) logError(err error) {
	fmt.Fprintf(l.output, "! ERROR: %v", err)
}
