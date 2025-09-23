package dto

type PostReq struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}
type CommentReq struct {
	Content string `json:"content"`
}

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
}
