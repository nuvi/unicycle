package unicycle

import (
	"testing"
)

func TestTrim(t *testing.T) {
	trimmed := Trim(make([]string, 10, 100))
	if len(trimmed) != cap(trimmed) {
		t.Errorf("len(trimmed) != cap(trimmed)")
	}
}
