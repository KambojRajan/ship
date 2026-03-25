package trace

import (
	"encoding/json"
	"io"
	"sync"
)

type jsonRecord struct {
	Name       string  `json:"name"`
	DurationMs float64 `json:"duration_ms"`
	Status     string  `json:"status"`
	Error      string  `json:"error,omitempty"`
}

type JSONSink struct {
	enc *json.Encoder
	mu  sync.Mutex
}

func NewJSONSink(w io.Writer) *JSONSink {
	return &JSONSink{enc: json.NewEncoder(w)}
}

func (s *JSONSink) Emit(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := jsonRecord{
		Name:       e.Name,
		DurationMs: float64(e.Duration.Microseconds()) / 1000.0,
		Status:     "ok",
	}
	if e.Err != nil {
		rec.Status = "error"
		rec.Error = e.Err.Error()
	}
	_ = s.enc.Encode(rec)
}
