package models

import (
	"time"
	"gorm.io/gorm"
)

type Folder struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	ParentID  *uint          `json:"parent_id"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	IconType  string         `json:"icon_type" gorm:"default:folder"`
	Path      string         `json:"path" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	User       User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent     *Folder  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Subfolders []Folder `json:"subfolders,omitempty" gorm:"foreignKey:ParentID"`
	Files      []File   `json:"files,omitempty" gorm:"foreignKey:FolderID"`
}