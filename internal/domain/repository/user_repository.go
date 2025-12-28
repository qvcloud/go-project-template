package repository

import (
	"github.com/qvcloud/go-project-template/internal/domain/entity"
)

type UserRepository interface {
	Query(id int64) *entity.User
}
