// Package cli implements a colored text handler suitable for command-line interfaces.
package loghandler

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/apex/log"
)

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// Default handler outputting to stderr.
var Default = New(os.Stderr)

type field struct {
	Name  string
	Value interface{}
}

type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	level := Strings[e.Level]
	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "%s %s %d %s", time.Now().Format("2006-01-02 15:04:05.000"), level, e.Level, e.Message)

	for _, f := range fields {
		fmt.Fprintf(h.Writer, " %s=\"%v\"", f.Name, f.Value)
	}

	fmt.Fprintln(h.Writer)

	return nil
}
