package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username string `gorm:"index:,unique"`
	Password string
}

func (u UserModel) TableName() string {
	return "user"
}

type UserModelSerializer struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

func (u UserModel) Serializer() UserModelSerializer {
	return UserModelSerializer{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Truncate(time.Second),
		UpdatedAt: u.UpdatedAt.Truncate(time.Second),
		Username:  u.Username,
	}
}
