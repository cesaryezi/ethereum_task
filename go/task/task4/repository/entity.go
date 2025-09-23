package repository

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Email    string    `gorm:"unique;not null"`
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	PostID  uint
	Post    Post `gorm:"foreignKey:PostID"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
}
