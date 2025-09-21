package task

import "testing"

func TestFindOneElement(t *testing.T) {
	nums := []int{4, 1, 2, 1, 2}
	res := findOneElement(nums)
	t.Log(res)
}

func TestIsPalindrome(t *testing.T) {
	nums := 12321
	res := isPalindrome(nums)
	t.Log(res)
}

func TestIsValid(t *testing.T) {
	nums := "([]{[()]][})"
	res := isValid(nums)
	t.Log(res)
}

func TestLongestCommonPrefix(t *testing.T) {
	strs := []string{"flower", "flow", "flight"}
	res := longestCommonPrefix(strs)
	t.Log(res)
}

func TestPlusOne(t *testing.T) {
	strs := []int{1, 2, 3}
	res := plusOne(strs)
	t.Log(res)
}

func TestRemoveDuplicates(t *testing.T) {
	strs := []int{1, 2, 3, 3, 5, 6, 7, 7, 8, 8, 8}
	res := removeDuplicates(strs)
	t.Log(res)
}

func TestTwoSum(t *testing.T) {
	strs := []int{1, 2, 3, 4, 5, 6, 9, 11}
	res := twoSum(strs, 20)
	t.Log(res)
}
