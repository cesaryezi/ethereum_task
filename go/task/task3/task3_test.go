package task3

import "testing"

func TestInsertData(t *testing.T) {
	insertData()

}

func TestFindUserA(t *testing.T) {
	user := findUserA(4)
	user.Info()
}

func TestFindMostCommentsPost(t *testing.T) {
	post, _ := findMostCommentsPost()
	post.Info()
}
