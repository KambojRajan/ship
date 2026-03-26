// Package trace provides lightweight step-level tracing for Ship commands.
//
// A global Sink (default: no-op) receives an Event each time a named step
// completes. cmd/trace.go swaps in a real Sink via SetSink before running a
// command, then restores the no-op when done.
package trace

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Event is emitted by a pipeline step when it completes.
type Event struct {
	Name     string
	Duration time.Duration
	Err      error
}

// Sink receives trace events.
type Sink interface {
	Emit(Event)
}

// ---- global sink --------------------------------------------------------

var (
	mu      sync.Mutex
	current Sink = noopSink{}
)

type noopSink struct{}

func (noopSink) Emit(Event) {}

// SetSink replaces the active sink and returns a function that restores the
// previous one.  Call it with defer in trace-aware callers.
func SetSink(s Sink) func() {
	mu.Lock()
	prev := current
	current = s
	mu.Unlock()
	return func() {
		mu.Lock()
		current = prev
		mu.Unlock()
	}
}

// ---- step helpers -------------------------------------------------------

// Step records the start time of a named pipeline step and returns a
// completion function.  Call the completion function with the step's final
// error (nil on success) to emit the event.
//
//	end := trace.Step("LoadIndex")
//	index, err := entities.LoadIndex(path)
//	end(err)
func Step(name string) func(error) {
	start := time.Now()
	return func(err error) {
		mu.Lock()
		s := current
		mu.Unlock()
		s.Emit(Event{
			Name:     name,
			Duration: time.Since(start),
			Err:      err,
		})
	}
}

// Meta records a key/value annotation as a zero-duration event.
// It is used to attach context (e.g. file counts) to a trace.
func Meta(key, value string) {
	mu.Lock()
	s := current
	mu.Unlock()
	s.Emit(Event{
		Name: fmt.Sprintf("meta:%s=%s", key, value),
	})
}

// ---- PrettySink ---------------------------------------------------------

// PrettySink prints coloured step output to a writer and accumulates a
// summary that can be printed after all steps complete.
type PrettySink struct {
	mu      sync.Mutex
	w       io.Writer
	steps   []stepRecord
	metas   []string
}

type stepRecord struct {
	name string
	dur  time.Duration
	err  error
}

func NewPrettySink(w io.Writer) *PrettySink {
	return &PrettySink{w: w}
}

func (p *PrettySink) Emit(e Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if e.Duration == 0 {
		// meta annotation
		p.metas = append(p.metas, e.Name)
		return
	}

	p.steps = append(p.steps, stepRecord{e.Name, e.Duration, e.Err})

	status := "\033[32m✓\033[0m"
	if e.Err != nil {
		status = "\033[31m✗\033[0m"
	}
	fmt.Fprintf(p.w, "  %s  %-36s  %s\n",
		status,
		e.Name,
		formatDuration(e.Duration),
	)
}

// PrintSummary prints a footer with total duration and any recorded metadata.
func (p *PrettySink) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var total time.Duration
	errors := 0
	for _, s := range p.steps {
		total += s.dur
		if s.err != nil {
			errors++
		}
	}

	fmt.Fprintf(p.w, "\n  \033[90m%d steps · %s", len(p.steps), formatDuration(total))
	if errors > 0 {
		fmt.Fprintf(p.w, " · \033[31m%d error(s)\033[90m", errors)
	}
	if len(p.metas) > 0 {
		fmt.Fprintf(p.w, " · ")
		for i, m := range p.metas {
			if i > 0 {
				fmt.Fprint(p.w, ", ")
			}
			fmt.Fprint(p.w, m)
		}
	}
	fmt.Fprintln(p.w, "\033[0m\n")
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000.0)
}
