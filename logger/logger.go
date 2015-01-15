package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger interface {
	Log(fields map[string]interface{})
	LogFromClientRequest(fields map[string]interface{}, req *http.Request)
}

type logEntry struct {
	Timestamp time.Time              `json:"@timestamp"`
	Fields    map[string]interface{} `json:"@fields"`
}

type jsonLogger struct {
	writer io.Writer
	lines  chan *[]byte
}

// New creates a new Logger.   The output variable sets the
// destination to which log data will be written.  This can be
// either an io.Writer, or a string.  With the latter, this is either
// one of "STDOUT" or "STDERR", or the path to the file to log to.
func New(output interface{}) (logger Logger, err error) {
	l := &jsonLogger{}
	l.writer, err = openWriter(output)
	if err != nil {
		return nil, err
	}
	l.lines = make(chan *[]byte, 100)
	go l.writeLoop()
	return l, nil
}

func openWriter(output interface{}) (w io.Writer, err error) {
	switch out := output.(type) {
	case io.Writer:
		w = out
	case string:
		if out == "STDERR" {
			w = os.Stderr
		} else if out == "STDOUT" {
			w = os.Stdout
		} else {
			w, err = os.OpenFile(out, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("Invalid output type %T(%v)", output, output)
	}
	return
}

func (l *jsonLogger) writeLoop() {
	for {
		line := <-l.lines
		_, err := l.writer.Write(*line)
		if err != nil {
			log.Printf("entity-extractor: Error writing to error log: %v", err)
		}
	}
}

func (l *jsonLogger) writeLine(line []byte) {
	line = append(line, 10) // Append a newline
	l.lines <- &line
}

func (l *jsonLogger) Log(fields map[string]interface{}) {
	entry := &logEntry{time.Now(), fields}
	line, err := json.Marshal(entry)
	if err != nil {
		log.Printf("entity-extractor/logger: Error encoding JSON: %v", err)
	}
	l.writeLine(line)
}

func (l *jsonLogger) LogFromClientRequest(fields map[string]interface{}, req *http.Request) {
	fields["request_method"] = req.Method
	fields["request"] = fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, req.Proto)
	fields["varnish_id"] = req.Header.Get("X-Varnish")

	l.Log(fields)
}
