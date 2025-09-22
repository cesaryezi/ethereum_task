package task2

import "testing"

func TestPointerUse(t *testing.T) {
	num := 12
	pointerUse(&num)
	t.Log(num)
}

func TestSliceMul(t *testing.T) {
	num := []int{1, 2, 3, 4, 5}
	sliceMul(&num)
	t.Log(num)
}

func TestGoroutine(t *testing.T) {
	goroutine()
}
