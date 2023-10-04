package unicycle

import (
	"testing"
)

func TestSum(t *testing.T) {
	if total := Sum([]int{1, 2, 3, 4, 5}); total != 15 {
		t.Errorf("Sum() returned wrong total %v", total)
	}
	if total := Sum([]float32{1.0, 3.4, 5.6, 7.9, 9.1}); total != 27.0 {
		t.Errorf("Sum() returned wrong total %v", total)
	}
}