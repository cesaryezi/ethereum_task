package task

import "strings"

/*
136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，
其余每个元素均出现两次。找出那个只出现了一次的元素。可以使用 for 循环遍历数组，
结合 if 条件判断和 map 数据结构来解决，
例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
*/
func findOneElement(arr []int) (res int) {
	var resMap = make(map[int]int)
	for _, num := range arr {
		resMap[num]++
	}
	for k, v := range resMap {
		if v == 1 {
			res = k
		}
	}
	return
}

//回文数
//
//考察：数字操作、条件判断
//题目：判断一个整数是否是回文数

func isPalindrome(x int) bool {

	if x < 0 {
		return false
	}

	var res int
	original := x

	for x > 0 {
		res = res*10 + x%10
		x /= 10
	}

	return res == original
}

// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
//
// 有效字符串需满足：
//
// 左括号必须用相同类型的右括号闭合。
// 左括号必须以正确的顺序闭合。
// 每个右括号都有一个对应的相同类型的左括号。
func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}

	var stack []int32
	pairs := map[int32]int32{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		if pair, ok := pairs[char]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != pair {
				return false
			}
			//弹出左括号
			stack = stack[:len(stack)-1]
		} else {
			// 左括号加入
			stack = append(stack, char)
		}
	}

	return len(stack) == 0
}

// 查找字符串数组中的最长公共前缀
func longestCommonPrefix(s []string) string {

	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}

	res := s[0]
	for i := 1; i < len(s); i++ {
		for !strings.HasPrefix(s[i], res) {
			res = res[:len(res)-1]
			if len(res) == 0 {
				return ""
			}
		}

	}
	return res
}

func plusOne(digits []int) []int {

	if len(digits) == 0 {
		return nil
	}

	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		} else {
			digits[i] = 0
		}
	}
	digits = append([]int{1}, digits...)
	return digits

}

func removeDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}
	i := 0
	for j := 1; j < len(nums)-1; j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]

		}
	}
	return i + 1

}

func twoSum(nums []int, target int) []int {
	if len(nums) == 0 {
		return nil
	}
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
