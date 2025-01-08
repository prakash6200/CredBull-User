package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ReferralCode string        `gorm:"unique;not null"` // For referral
	ProfileImage string        `gorm:"default:''"`
	Name         string        `gorm:"default:''"`
	Email        string        `gorm:"unique;not null"`
	Mobile       string        `gorm:"default:''"`
	Role         string        `gorm:"default:'savings'"`
	Password     string        `gorm:"not null"`
	BankDetails  []BankDetails `gorm:"foreignKey:UserID"` // Corrected foreign key reference
	UserKYC      UserKYC       `gorm:"foreignKey:UserID"` // Corrected foreign key reference
	IsDeleted    bool          `gorm:"default:false"`
}
