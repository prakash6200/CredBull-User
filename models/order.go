package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model           // Embeds ID, CreatedAt, UpdatedAt, DeletedAt fields
	ProductRefer int     `json:"product_id"`
	Product      Product `gorm:"foreignKey:ProductRefer"`
	UserRefer    int     `json:"user_id"`
	User         User    `gorm:"foreignKey:UserRefer"`
}
