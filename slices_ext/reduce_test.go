package slices_ext

import (
	"testing"
)

func sum(total, current int) int {
	return total + current
}

func TestReduce(t *testing.T) {
	result := Reduce([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, sum, 0)
	if result != 45 {
		t.Errorf("Reduce() returned unexpected %d", result)
	}

	if Reduce(nil, sum, 7) != 7 {
		t.Error("Reduce(nil, reducer, initial) should return the initial value")
	}
}
