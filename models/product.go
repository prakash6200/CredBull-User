package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model          // Embeds ID, CreatedAt, UpdatedAt, DeletedAt fields
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}
