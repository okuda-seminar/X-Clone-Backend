package entities

import "github.com/google/uuid"

type Like struct {
	PostID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	UserID string    `gorm:"type:varchar(255);not null;primaryKey"`
}
