package usecase

import (
	"github.com/qvcloud/go-project-template/internal/domain/entity"
	"github.com/qvcloud/go-project-template/internal/domain/repository"
	"go.uber.org/zap"
)

type UserUseCase interface {
	GetUser(id int64) *entity.User
}

type userUseCase struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserUseCase(repo repository.UserRepository, logger *zap.Logger) UserUseCase {
	return &userUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (u *userUseCase) GetUser(id int64) *entity.User {
	return u.repo.Query(id)
}
