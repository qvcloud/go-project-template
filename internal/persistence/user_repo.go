package persistence

import (
	"github.com/qvcloud/go-project-template/internal/domain/entity"
	"github.com/qvcloud/go-project-template/internal/domain/repository"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &user{
		db: db,
	}
}

func (u *user) Query(_ int64) *entity.User {
	return &entity.User{}
}
