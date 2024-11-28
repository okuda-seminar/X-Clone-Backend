package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents an entry of `users` table.
// It contains properties such as Username and DisplayName,
// which are associated with a user's identity
// Username is a unique identifier for a user's account and is preceded by the "@" symbol.
// It's used in the URL of the user's profile
// and is how other users mention or reference them in tweets.
//
// DisplayName is the name that appears on a user's profile and alongside their tweets.
// It doesn't need to be unique and can be changed by the user.
//
// For more information on terminology, refer to: https://help.twitter.com/en/resources/glossary.
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username    string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	DisplayName string    `gorm:"type:varchar(255);not null" json:"display_name"`
	Bio         string    `gorm:"type:text" json:"bio"`
	IsPrivate   bool      `gorm:"type:boolean;default:false" json:"is_private"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
}
