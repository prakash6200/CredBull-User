package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt fields
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"password"`
}
