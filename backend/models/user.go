package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GoogleID  string         `gorm:"type:varchar(255);uniqueIndex" json:"google_id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	PhotoURL  string         `gorm:"type:text" json:"photo_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
