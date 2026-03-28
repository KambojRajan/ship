package trace

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewPrettySink(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
		want  *PrettySink
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got := NewPrettySink(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewPrettySink() gotW = %v, want %v", gotW, tt.wantW)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrettySink() = %v, want %v", got, tt.want)
			}
		})
	}
}
