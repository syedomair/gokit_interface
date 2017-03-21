package models

import "time"

type Book struct {
	Id          int64     `json:"id" gorm:"column:id"`
	UserId      int64     `json:"user_id" gorm:"column:user_id"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Publish     bool      `json:"publish" gorm:"column:publish"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Book) TableName() string {
	return "book"
}
