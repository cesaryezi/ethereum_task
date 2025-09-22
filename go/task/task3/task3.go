package task3

import (
	"errors"
	"fmt"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

//Sqlx入门
//题目1：使用SQL扩展库进行查询
//假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
//要求 ：
//编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。

//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

//题目2：实现类型安全映射
//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

// 进阶gorm
// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
// Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
type User struct {
	gorm.Model
	Name       string
	Posts      []Post    `gorm:"foreignKey:UserID"`
	Comments   []Comment `gorm:"foreignKey:UserID"`
	PostsCount int
}

func (u *User) Info() {

	fmt.Println("User:", u.Name)
	fmt.Println("Posts:")
	for _, post := range u.Posts {
		fmt.Println("-", post.Title)
		fmt.Println("  Comments:")
		for _, comment := range post.Comments {
			fmt.Println("  -", comment.Content)
		}
	}
	fmt.Println("Comments:")
	for _, comment := range u.Comments {
		fmt.Println("-", comment.Content)
	}

}

type Post struct {
	gorm.Model
	Title         string
	Content       string
	UserID        uint
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentsCount int
	Status        string
}

func (post *Post) Info() {
	fmt.Println("Post:", post.Title)
	fmt.Println("  Comments:")
	for _, comment := range post.Comments {
		fmt.Println("  -", comment.Content)
	}
	fmt.Println("User:", post.User.Name)
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	Post    Post `gorm:"foreignKey:PostID"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
}

// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func insertData() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database instance")
		}
		sqlDB.Close()
	}()

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic("failed to migrate database schema")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建用户
		user := &User{Name: "jinzhu"}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. 创建文章
		posts := []Post{
			{Title: "Post 1", Content: "Content 1", UserID: user.ID},
			{Title: "Post 2", Content: "Content 2", UserID: user.ID},
		}
		if err := tx.Create(&posts).Error; err != nil {
			return err
		}

		// 3. 创建评论并关联到文章
		comments := []Comment{
			{Content: "Comment 1 on Post 1", PostID: posts[0].ID, UserID: user.ID},
			{Content: "Comment 2 on Post 1", PostID: posts[0].ID, UserID: user.ID},
			{Content: "Comment 1 on Post 2", PostID: posts[1].ID, UserID: user.ID},
		}
		return tx.Create(&comments).Error
	})

	if err != nil {
		panic("failed to create user data")
	}

}

// 查询某个用户发布的所有文章及其对应的评论信息
func findUserA(userId int) *User {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database instance")
		}
		sqlDB.Close()
	}()

	var user User
	db.Preload("Posts").Preload("Posts.Comments").First(&user, userId)

	return &user

}

// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
func findMostCommentsPost() (*Post, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database instance")
		}
		sqlDB.Close()
	}()

	//select posts.* ,count(comments.id) as commentCnt from posts
	//left join comments on posts.id = comments.post_id
	//group by posts.id
	//order by commentCnt desc
	var post Post
	err2 := db.Model(&Post{}).
		Select("posts.*,count(comments.id) as commentCnt").
		Joins("left join comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("commentCnt desc").
		Limit(1).
		Preload("Comments").
		Preload("User").
		First(&post).
		Error

	return &post, err2
}

// 题目3：钩子函数
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，
// 如果评论数量为 0， 则更新文章的评论状态为 "无评论"。
func (u *Post) AfterCreate(tx *gorm.DB) (err error) {
	var user User
	tx.First(&user, u.UserID)
	user.PostsCount++
	db := tx.Model(&User{}).Where("id = ?", u.UserID).Update("posts_count", user.PostsCount)

	if db.Error != nil {
		return errors.New("rollback invalid op")
	}
	return nil

}

func (comment *Comment) AfterDelete(tx *gorm.DB) (err error) {

	var post Post
	tx.First(&post, comment.PostID)

	var db *gorm.DB
	if post.CommentsCount > 0 {
		post.CommentsCount--
		db = tx.Model(&Post{}).Where("id = ?", comment.PostID).Updates(map[string]interface{}{"comments_count": post.CommentsCount, "status": ""})
	} else {
		if post.Status != "无评论" {
			db = tx.Model(&Post{}).Where("id = ?", comment.PostID).Update("status", "无评论")
		}
	}

	if db.Error != nil {
		return errors.New("rollback invalid op")
	}

	return nil

}

func (comment *Comment) AfterCreate(tx *gorm.DB) (err error) {

	var post Post
	tx.First(&post, comment.PostID)
	post.CommentsCount++
	db := tx.Model(&Post{}).Where("id = ?", comment.PostID).Update("comments_count", post.CommentsCount)
	if db.Error != nil {
		return errors.New("rollback invalid op")
	}

	return nil

}
