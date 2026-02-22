package usecase

import (
	"Assignment2/pkg/modules"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	createdUser   *modules.User
	createErr     error
	returnID      int
	getUser       *modules.User
	getUserErr    error
	updatedID     int
	updatedUser   *modules.User
	updateErr     error
	deletedID     int
	deleteRows    int64
	deleteErr     error
	createAudit   string
	createAuditID int
}

func (m *mockRepo) GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error) {
	return nil, nil
}
func (m *mockRepo) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.getUser, nil
}
func (m *mockRepo) CreateUser(ctx context.Context, u *modules.User) (int, error) {
	return 0, errors.New("not used")
}
func (m *mockRepo) UpdateUser(ctx context.Context, id int, u *modules.User) error {
	m.updatedID = id
	m.updatedUser = u
	return m.updateErr
}
func (m *mockRepo) DeleteUser(ctx context.Context, id int) (int64, error) {
	m.deletedID = id
	return m.deleteRows, m.deleteErr
}
func (m *mockRepo) CreateUserWithAudit(ctx context.Context, u *modules.User, action string) (int, error) {
	m.createdUser = u
	m.createAudit = action
	if m.createErr != nil {
		return 0, m.createErr
	}
	if m.createAuditID != 0 {
		return m.createAuditID, nil
	}
	return 1, nil
}

func TestCreateUser_HashesPasswordAndCallsRepo(t *testing.T) {
	m := &mockRepo{createAuditID: 7}
	uc := NewUserUsecase(m)

	id, err := uc.CreateUser(context.Background(), "John", nil, nil, "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 7 {
		t.Fatalf("unexpected id: %d", id)
	}
	if m.createdUser == nil {
		t.Fatal("repo user is nil")
	}
	if m.createdUser.PasswordHash == "" {
		t.Fatal("password hash is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(m.createdUser.PasswordHash), []byte("secret123")); err != nil {
		t.Fatalf("password hash mismatch: %v", err)
	}
	if m.createAudit != "user_created" {
		t.Fatalf("unexpected audit action: %s", m.createAudit)
	}
}

func TestCreateUser_InvalidInput(t *testing.T) {
	m := &mockRepo{}
	uc := NewUserUsecase(m)

	_, err := uc.CreateUser(context.Background(), "", nil, nil, "")
	if !errors.Is(err, modules.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestAuthLogin_Success(t *testing.T) {
	h, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	m := &mockRepo{
		getUser: &modules.User{ID: 2, PasswordHash: string(h)},
	}
	a := NewAuthUsecase(m, "jwt-secret", time.Minute)

	token, err := a.Login(context.Background(), 2, "pass123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestAuthLogin_Unauthorized(t *testing.T) {
	h, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	m := &mockRepo{
		getUser: &modules.User{ID: 2, PasswordHash: string(h)},
	}
	a := NewAuthUsecase(m, "jwt-secret", time.Minute)

	_, err := a.Login(context.Background(), 2, "bad")
	if !errors.Is(err, modules.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}
