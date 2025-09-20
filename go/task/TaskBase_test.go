package task

import "testing"

func TestFindOneElement(t *testing.T) {
	nums := []int{4, 1, 2, 1, 2}
	res := findOneElement(nums)
	t.Log(res)
}
