package entity

import (
	"time"

	"github.com/google/uuid"
)

type UUID struct {
	ID string `gorm:"primarykey;type:varchar(36)" json:"id"`
}

type ID struct {
	ID int64 `gorm:"primarykey,autoIncrement" json:"id"`
}

func (u *UUID) GenerateID() {
	if len(u.ID) == 0 {
		u.ID = uuid.NewString()
	}
}

type BaseAt struct {
	CreatedAt time.Time `json:"createdAt" gorm:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"updatedAt"`
}
