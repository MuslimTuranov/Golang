package users

import (
	pg "Assignment2/internal/repository/postgres"
	"Assignment2/pkg/modules"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type UserRepo struct {
	db *pg.Dialect
}

func NewUserRepo(db *pg.Dialect) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUsers(ctx context.Context, limit, offset int) ([]modules.User, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var users []modules.User
	q := `SELECT id, name, email, age, created_at, updated_at, deleted_at
	      FROM users
	      WHERE deleted_at IS NULL
	      ORDER BY id
	      LIMIT $1 OFFSET $2`
	if err := r.db.DB.SelectContext(ctx, &users, q, limit, offset); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	var u modules.User
	q := `SELECT id, name, email, age, password_hash, created_at, updated_at, deleted_at
	      FROM users
	      WHERE id=$1 AND deleted_at IS NULL`
	err := r.db.DB.GetContext(ctx, &u, q, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: user id=%d", modules.ErrNotFound, id)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, u *modules.User) (int, error) {
	if u == nil || u.Name == "" || u.PasswordHash == "" {
		return 0, modules.ErrInvalidInput
	}

	q := `INSERT INTO users (name, email, age, password_hash, created_at, updated_at)
	      VALUES ($1, $2, $3, $4, now(), now())
	      RETURNING id`

	var id int
	if err := r.db.DB.QueryRowxContext(ctx, q, u.Name, u.Email, u.Age, u.PasswordHash).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, id int, u *modules.User) error {
	if u == nil {
		return modules.ErrInvalidInput
	}

	q := `UPDATE users
	      SET name = COALESCE($2, name),
	          email = COALESCE($3, email),
	          age = COALESCE($4, age),
	          password_hash = COALESCE($5, password_hash),
	          updated_at = now()
	      WHERE id=$1 AND deleted_at IS NULL`
	namePtr := sql.NullString{}
	if u.Name != "" {
		namePtr = sql.NullString{String: u.Name, Valid: true}
	}

	res, err := r.db.DB.ExecContext(ctx, q, id,
		nullString(namePtr),
		nullStringPtr(u.Email),
		nullIntPtr(u.Age),
		nullPassword(u.PasswordHash),
	)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return fmt.Errorf("%w: user id=%d", modules.ErrNotFound, id)
	}
	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int) (int64, error) {
	q := `UPDATE users SET deleted_at=now(), updated_at=now()
	      WHERE id=$1 AND deleted_at IS NULL`
	res, err := r.db.DB.ExecContext(ctx, q, id)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if ra == 0 {
		return 0, fmt.Errorf("%w: user id=%d", modules.ErrNotFound, id)
	}
	return ra, nil
}

func (r *UserRepo) CreateUserWithAudit(ctx context.Context, u *modules.User, action string) (int, error) {
	tx, err := r.db.DB.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var id int
	insUser := `INSERT INTO users (name, email, age, password_hash, created_at, updated_at)
	            VALUES ($1, $2, $3, $4, now(), now())
	            RETURNING id`
	if err := tx.QueryRowxContext(ctx, insUser, u.Name, u.Email, u.Age, u.PasswordHash).Scan(&id); err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(map[string]any{
		"name":  u.Name,
		"email": u.Email,
		"age":   u.Age,
	})
	insAudit := `INSERT INTO audit_logs(action, user_id, payload, created_at)
	             VALUES ($1, $2, $3, now())`
	if _, err := tx.ExecContext(ctx, insAudit, action, id, payload); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func nullString(ns sql.NullString) any {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func nullStringPtr(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullIntPtr(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullPassword(v string) any {
	if v == "" {
		return nil
	}
	return v
}
