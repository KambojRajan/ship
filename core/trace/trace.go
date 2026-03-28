// Package trace provides lightweight step-level tracing for Ship commands.
//
// A global Sink (default: no-op) receives an Event each time a named step
// completes. cmd/trace.go swaps in a real Sink via SetSink before running a
// command, then restores the no-op when done.
package trace

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ANSI color codes used by PrettySink.
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorCyan  = "\033[36m"
	colorGray  = "\033[90m"
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
// previous one. Call it with defer in trace-aware callers.
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

// activeSink returns the current sink under the lock.
func activeSink() Sink {
	mu.Lock()
	s := current
	mu.Unlock()
	return s
}

// ---- step helpers -------------------------------------------------------

// Step records the start time of a named pipeline step and returns a
// completion function. Call the completion function with the step's final
// error (nil on success) to emit the event.
//
//	end := trace.Step("LoadIndex")
//	index, err := entities.LoadIndex(path)
//	end(err)
func Step(name string) func(error) {
	start := time.Now()
	return func(err error) {
		activeSink().Emit(Event{
			Name:     name,
			Duration: time.Since(start),
			Err:      err,
		})
	}
}

// Meta records a key/value annotation as a zero-duration event.
// It is used to attach context (e.g. file counts) to a trace.
func Meta(key, value string) {
	activeSink().Emit(Event{Name: fmt.Sprintf("meta:%s=%s", key, value)})
}

// ---- PrettySink ---------------------------------------------------------

// PrettySink accumulates step events and renders a coloured summary table
// when PrintSummary is called.
type PrettySink struct {
	mu    sync.Mutex
	w     io.Writer
	steps []stepRecord
	metas []string
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
		p.metas = append(p.metas, e.Name)
		return
	}
	p.steps = append(p.steps, stepRecord{e.Name, e.Duration, e.Err})
}

// PrintSummary prints all steps with a relative time bar and a footer summary.
func (p *PrettySink) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	total, errCount := p.totals()

	const barWidth = 18
	for _, s := range p.steps {
		status := colorGreen + "✓" + colorReset
		if s.err != nil {
			status = colorRed + "✗" + colorReset
		}

		var frac float64
		if total > 0 {
			frac = float64(s.dur) / float64(total)
		}
		bar, pct := buildBar(frac, barWidth)

		fmt.Fprintf(p.w, "  %s  %-36s  %8s  %s  %s%3.0f%%%s\n",
			status, s.name, formatDuration(s.dur), bar, colorGray, pct, colorReset)
	}

	fmt.Fprintf(p.w, "\n  %s%d steps · %s", colorGray, len(p.steps), formatDuration(total))
	if errCount > 0 {
		fmt.Fprintf(p.w, " · %s%d error(s)%s", colorRed, errCount, colorGray)
	}
	fmt.Fprintln(p.w, colorReset)

	if len(p.metas) > 0 {
		fmt.Fprint(p.w, "  "+colorGray)
		for i, m := range p.metas {
			if i > 0 {
				fmt.Fprint(p.w, "  ")
			}
			fmt.Fprint(p.w, strings.TrimPrefix(m, "meta:"))
		}
		fmt.Fprintln(p.w, colorReset)
	}
	fmt.Fprintln(p.w)
}

// totals returns the sum of all step durations and the number of failed steps.
func (p *PrettySink) totals() (total time.Duration, errCount int) {
	for _, s := range p.steps {
		total += s.dur
		if s.err != nil {
			errCount++
		}
	}
	return
}

// buildBar renders a coloured progress bar of the given width for frac ∈ [0,1].
// It returns the bar string and the percentage value.
func buildBar(frac float64, width int) (bar string, pct float64) {
	pct = frac * 100
	filled := int(frac*float64(width) + 0.5)
	if filled > width {
		filled = width
	}
	bar = colorCyan + strings.Repeat("█", filled) +
		colorGray + strings.Repeat("░", width-filled) + colorReset
	return
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000.0)
}
