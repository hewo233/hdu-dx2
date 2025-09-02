package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string       `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password string       `json:"-" gorm:"size:100;not null"`
	Phone    string       `json:"phone" gorm:"size:11;not null"`
	Families []FamilyUser `json:"families" gorm:"foreignKey:UserID"`
}

func NewUser() *User {
	return &User{}
}

type Family struct {
	gorm.Model
	Name     string       `json:"name" gorm:"size:100;not null"`
	Users    []FamilyUser `json:"users" gorm:"foreignKey:FamilyID"`
	Password string       `json:"-" gorm:"size:100;not null"` // for joining family
}

func NewFamily() *Family {
	return &Family{}
}

type FamilyUser struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey"`
	FamilyID uint   `json:"family_id" gorm:"primaryKey"`
	Role     string `json:"role" gorm:"size:20;not null"` // father, mother, son, daughter...

	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt gorm.DeletedAt `gorm:"index"`
}
