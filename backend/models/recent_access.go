package models

import (
	"time"
	"gorm.io/gorm"
)

type RecentAccess struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	ItemType   string    `json:"item_type" gorm:"not null;check:item_type IN ('file','folder')"`
	ItemID     uint      `json:"item_id" gorm:"not null"`
	AccessedAt time.Time `json:"accessed_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Add unique constraint to prevent duplicate recent access entries
func (RecentAccess) TableName() string {
	return "recent_accesses"
}

// BeforeCreate hook to set accessed_at time
func (ra *RecentAccess) BeforeCreate(tx *gorm.DB) error {
	ra.AccessedAt = time.Now()
	return nil
}

// BeforeUpdate hook to update accessed_at time
func (ra *RecentAccess) BeforeUpdate(tx *gorm.DB) error {
	ra.AccessedAt = time.Now()
	return nil
}