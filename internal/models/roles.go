package models

import (
	"time"
)

type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:100;uniqueIndex;not null"` // role name like "admin", "user", etc.
	Description string `gorm:"size:255"`                      // optional description of the role
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"` // optional soft delete
}
