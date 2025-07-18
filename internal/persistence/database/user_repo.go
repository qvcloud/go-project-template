package database

import (
	"github.com/qvcloud/go-project-template/internal/persistence"
	"github.com/qvcloud/go-project-template/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) persistence.UserRepository {
	return &user{
		db: db,
	}
}

func (u *user) Query(id int64) *entities.User {
	return &entities.User{}
}
