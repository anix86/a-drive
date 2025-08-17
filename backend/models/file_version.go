package models

import (
	"time"
	"gorm.io/gorm"
)

type FileVersion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FileID      uint           `json:"file_id" gorm:"not null"`
	Version     int            `json:"version" gorm:"not null"`
	FilePath    string         `json:"file_path" gorm:"not null"`
	Size        int64          `json:"size" gorm:"not null"`
	Checksum    string         `json:"checksum"`
	Comment     string         `json:"comment"`
	CreatedBy   uint           `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	File      File `json:"file,omitempty" gorm:"foreignKey:FileID"`
	CreatedByUser User `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBy"`
}