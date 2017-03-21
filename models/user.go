package models

import "time"

type User struct {
	Id        int64     `json:"id" gorm:"column:id"`
	Email     string    `json:"email" gorm:"column:email"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Password  string    `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}
