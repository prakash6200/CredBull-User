package models

import (
	"gorm.io/gorm"
)

type UserKYC struct {
	gorm.Model
	UserID      uint     `gorm:"foreignKey:UserID"`               // Corrected foreign key reference
	Country     string   `gorm:"default:''"`                      // Country of the user
	AadharProof Aadhar   `gorm:"embedded;embeddedPrefix:aadhar_"` // Embedded struct for Aadhar details
	PanProof    PanProof `gorm:"embedded;embeddedPrefix:pan_"`    // Embedded struct for PAN details
	IsVerified  bool     `gorm:"default:false"`                   // KYC verification status
	IsDeleted   bool     `gorm:"default:false"`
}

type Aadhar struct {
	AadharNumber string `gorm:"unique;not null"` // Aadhar number must be unique and not null
	Name         string `gorm:"default:''"`      // Name on the Aadhar card
	DOB          string `gorm:"default:''"`      // Date of Birth
	Address      string `gorm:"default:''"`      // Address on the Aadhar card
	IsVerified   bool   `gorm:"default:false"`   // Verification status, default is false
	RefID        string `gorm:"default:''"`      // Reference ID
}

type PanProof struct {
	PanNumber  string `gorm:"unique;not null"` // PAN number must be unique and not null
	Name       string `gorm:"default:''"`      // Name on the PAN card
	IsVerified bool   `gorm:"default:false"`   // Verification status, default is false
}
