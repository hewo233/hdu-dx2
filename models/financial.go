package models

import (
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	gorm.Model
	Date        time.Time `json:"date" gorm:"not null"`
	Type        string    `json:"type" gorm:"size:100;not null;oneof:income,expense"`
	Amount      int       `json:"amount" gorm:"not null"` // 分
	Category    string    `json:"category" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:255"`
	Object      string    `json:"object" gorm:"size:100;not null"` // 谁给的/给谁的
	Username    string    `json:"username" gorm:"size:100;not null"`
	FamilyID    uint      `json:"family_id" gorm:"not null;index"`
}

func NewBill() *Bill {
	return &Bill{}
}
