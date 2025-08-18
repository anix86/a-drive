package models

import (
	"time"
	"gorm.io/gorm"
)

type File struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	OriginalName string         `json:"original_name" gorm:"not null"`
	FolderID     *uint          `json:"folder_id"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	FilePath     string         `json:"file_path" gorm:"not null"`
	Size         int64          `json:"size" gorm:"not null"`
	MimeType     string         `json:"mime_type"`
	CurrentVersion int          `json:"current_version" gorm:"default:1"`
	VersioningEnabled bool      `json:"versioning_enabled" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	IsFavorite   bool           `json:"is_favorite" gorm:"-"`
	
	User     User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Folder   *Folder       `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
	Versions []FileVersion `json:"versions,omitempty" gorm:"foreignKey:FileID"`
}