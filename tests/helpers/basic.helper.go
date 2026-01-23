package helpers

import (
	"os"
	"testing"
)

func AssertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("expected %s to exist", path)
		}
		t.Fatalf("error checking %s: %v", path, err)
	}
}

func AssertNil(err error) {
	if err != nil {
		panic("expected error")
	}
}

func AssertNotNil(err error) {
	if err == nil {
		panic(err)
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
