package models

import (
	"time"
	"gorm.io/gorm"
)

type Favorite struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	ItemType string    `json:"item_type" gorm:"not null;check:item_type IN ('file','folder')"`
	ItemID   uint      `json:"item_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Add unique constraint to prevent duplicate favorites
func (Favorite) TableName() string {
	return "favorites"
}

// BeforeCreate hook to add unique constraint
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	// This will be enforced by database unique constraint
	return nil
}