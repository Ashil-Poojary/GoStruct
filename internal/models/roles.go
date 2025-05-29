package models

import (
	"time"
)

type Role struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;uniqueIndex;not null"` // role name like "admin", "user", etc.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // optional soft delete
}
