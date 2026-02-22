package repository

import (
	"Assignment2/pkg/modules"
	"context"
)

type UserRepository interface {
	GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, u *modules.User) (int, error)
	UpdateUser(ctx context.Context, id int, u *modules.User) error
	DeleteUser(ctx context.Context, id int) (int64, error)
	CreateUserWithAudit(ctx context.Context, u *modules.User, action string) (int, error)
}

type Repositories struct {
	UserRepo UserRepository
}

func NewRepositories(userRepo UserRepository) *Repositories {
	return &Repositories{UserRepo: userRepo}
}
