package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model                 // Auto includes ID, CreatedAt, UpdatedAt, DeletedAt
	UserID       string        `gorm:"unique;not null"` // For referral
	ProfileImage string        `gorm:"default:''"`
	Name         string        `gorm:"default:''"`
	Email        string        `gorm:"unique;not null"`
	Mobile       string        `gorm:"default:''"`
	Role         string        `gorm:"type:text;default:'savings'"`
	Password     string        `gorm:"not null"`
	BankDetails  []BankDetails `gorm:"foreignKey:UserID"` // One-to-Many association
}
