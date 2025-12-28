package service

import (
	"github.com/qvcloud/go-project-template/internal/domain/entity"
	"github.com/qvcloud/go-project-template/internal/domain/repository"
	"go.uber.org/zap"
)

type UserService interface {
	GetUser(id int64) *entity.User
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (u *userService) GetUser(id int64) *entity.User {
	return u.repo.Query(id)
}
