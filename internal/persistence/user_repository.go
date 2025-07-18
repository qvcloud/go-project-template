package persistence

import (
	"github.com/qvcloud/go-project-template/internal/persistence/database/entities"
)

type UserRepository interface {
	Query(id int64) *entities.User
}
