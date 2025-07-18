package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUID struct {
	Id string `gorm:"primarykey;type:varchar(36)" json:"id"`
}

type ID struct {
	Id int64 `gorm:"primarykey,autoIncrement" json:"id"`
}

func (u *UUID) BeforeCreate(tx *gorm.DB) error {
	if len(u.Id) == 0 {
		u.Id = uuid.NewString()
	}
	return nil
}

type BaseAt struct {
	CreatedAt time.Time `json:"createdAt" gorm:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"updatedAt"`
}
