package models

import (
	"time"
	"gorm.io/gorm"
)

type FileShare struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FileID      uint           `json:"file_id" gorm:"not null"`
	File        File           `json:"file,omitempty" gorm:"foreignKey:FileID"`
	SharedBy    uint           `json:"shared_by" gorm:"not null"`
	SharedByUser User          `json:"shared_by_user,omitempty" gorm:"foreignKey:SharedBy"`
	ShareToken  string         `json:"share_token" gorm:"uniqueIndex;not null"`
	ShareType   string         `json:"share_type" gorm:"not null"` // "public", "password", "private"
	Password    string         `json:"-" gorm:"column:password"`   // Only for password-protected shares
	ExpiresAt   *time.Time     `json:"expires_at"`                 // Optional expiration
	DownloadCount int          `json:"download_count" gorm:"default:0"`
	MaxDownloads  *int         `json:"max_downloads"`              // Optional download limit
	AllowPreview bool          `json:"allow_preview" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ShareAccess struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ShareID     uint           `json:"share_id" gorm:"not null"`
	Share       FileShare      `json:"share,omitempty" gorm:"foreignKey:ShareID"`
	AccessedBy  string         `json:"accessed_by"` // IP address or identifier
	AccessedAt  time.Time      `json:"accessed_at"`
	UserAgent   string         `json:"user_agent"`
	Action      string         `json:"action"` // "view", "download"
	CreatedAt   time.Time      `json:"created_at"`
}