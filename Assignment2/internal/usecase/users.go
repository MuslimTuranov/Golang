package usecase

import (
	"Assignment2/internal/repository"
	"Assignment2/pkg/modules"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}

func (u *UserUsecase) GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(ctx, limit, offset)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *UserUsecase) CreateUser(ctx context.Context, name string, email *string, age *int, password string) (int, error) {
	if name == "" || password == "" {
		return 0, modules.ErrInvalidInput
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, modules.ErrInternal
	}

	user := &modules.User{
		Name:         name,
		Email:        email,
		Age:          age,
		PasswordHash: string(h),
	}

	return u.repo.CreateUserWithAudit(ctx, user, "user_created")
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id int, name *string, email *string, age *int, password *string) error {
	upd := &modules.User{}
	if name != nil {
		upd.Name = *name
	}
	upd.Email = email
	upd.Age = age

	if password != nil && *password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return modules.ErrInternal
		}
		upd.PasswordHash = string(h)
	}

	return u.repo.UpdateUser(ctx, id, upd)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int) (int64, error) {
	return u.repo.DeleteUser(ctx, id)
}
