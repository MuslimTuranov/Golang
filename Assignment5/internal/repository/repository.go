package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"users-service/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type UserFilter struct {
	ID        *int
	Name      *string
	Email     *string
	Gender    *string
	BirthDate *string
	Status    string
}

type UserListParams struct {
	Page     int
	PageSize int

	OrderBy  string
	OrderDir string

	Filters UserFilter
}

type CursorParams struct {
	Limit    int
	Cursor   *int
	OrderBy  string
	OrderDir string
	Status   string
}

var allowedOrderColumns = map[string]string{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"gender":     "gender",
	"birth_date": "birth_date",
}

var allowedOperators = map[string]bool{
	"=":     true,
	">":     true,
	"<":     true,
	"ILIKE": true,
}

func validateOrderBy(orderBy string) string {
	if col, ok := allowedOrderColumns[orderBy]; ok {
		return col
	}
	return "id"
}

func validateOrderDir(orderDir string) string {
	switch strings.ToUpper(orderDir) {
	case "ASC":
		return "ASC"
	case "DESC":
		return "DESC"
	default:
		return "ASC"
	}
}

func (r *Repository) GetPaginatedUsers(params UserListParams) (model.PaginatedResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 5
	}

	orderBy := validateOrderBy(params.OrderBy)
	orderDir := validateOrderDir(params.OrderDir)

	whereParts := []string{}
	args := []any{}
	argPos := 1

	switch params.Filters.Status {
	case "", "active":
		whereParts = append(whereParts, "deleted_at IS NULL")
	case "deleted":
		whereParts = append(whereParts, "deleted_at IS NOT NULL")
	case "all":
	default:
		return model.PaginatedResponse{}, errors.New("invalid status filter")
	}

	if params.Filters.ID != nil {
		whereParts = append(whereParts, fmt.Sprintf("id = $%d", argPos))
		args = append(args, *params.Filters.ID)
		argPos++
	}

	if params.Filters.Name != nil {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Filters.Name+"%")
		argPos++
	}

	if params.Filters.Email != nil {
		whereParts = append(whereParts, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+*params.Filters.Email+"%")
		argPos++
	}

	if params.Filters.Gender != nil {
		whereParts = append(whereParts, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, *params.Filters.Gender)
		argPos++
	}

	if params.Filters.BirthDate != nil {
		if _, err := time.Parse("2006-01-02", *params.Filters.BirthDate); err != nil {
			return model.PaginatedResponse{}, errors.New("invalid birth_date format, use YYYY-MM-DD")
		}

		whereParts = append(whereParts, fmt.Sprintf("birth_date = $%d", argPos))
		args = append(args, *params.Filters.BirthDate)
		argPos++
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users %s`, whereSQL)
	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return model.PaginatedResponse{}, err
	}

	offset := (params.Page - 1) * params.PageSize
	listQuery := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date, deleted_at
		FROM users
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, orderDir, argPos, argPos+1)

	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return model.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &u.DeletedAt); err != nil {
			return model.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return model.PaginatedResponse{}, err
	}

	return model.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(user1ID, user2ID int) ([]model.User, error) {
	if user1ID == user2ID {
		return nil, errors.New("user1_id and user2_id must be different")
	}

	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date, u.deleted_at
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1
		  AND uf2.user_id = $2
		  AND u.deleted_at IS NULL
		ORDER BY u.id ASC
	`

	rows, err := r.db.Query(query, user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &u.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *Repository) SoftDeleteUser(id int) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetUsersByCursor(params CursorParams) (model.CursorPaginatedResponse, error) {
	if params.Limit <= 0 {
		params.Limit = 5
	}

	orderBy := validateOrderBy(params.OrderBy)
	if orderBy != "id" {
		orderBy = "id"
	}
	orderDir := validateOrderDir(params.OrderDir)

	whereParts := []string{}
	args := []any{}
	argPos := 1

	switch params.Status {
	case "", "active":
		whereParts = append(whereParts, "deleted_at IS NULL")
	case "deleted":
		whereParts = append(whereParts, "deleted_at IS NOT NULL")
	case "all":
	default:
		return model.CursorPaginatedResponse{}, errors.New("invalid status filter")
	}

	if params.Cursor != nil {
		if orderDir == "ASC" {
			whereParts = append(whereParts, fmt.Sprintf("id > $%d", argPos))
		} else {
			whereParts = append(whereParts, fmt.Sprintf("id < $%d", argPos))
		}
		args = append(args, *params.Cursor)
		argPos++
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date, deleted_at
		FROM users
		%s
		ORDER BY id %s
		LIMIT $%d
	`, whereSQL, orderDir, argPos)

	args = append(args, params.Limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return model.CursorPaginatedResponse{}, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &u.DeletedAt); err != nil {
			return model.CursorPaginatedResponse{}, err
		}
		users = append(users, u)
	}

	var nextCursor *int
	if len(users) > 0 {
		lastID := users[len(users)-1].ID
		nextCursor = &lastID
	}

	return model.CursorPaginatedResponse{
		Data:       users,
		NextCursor: nextCursor,
		Limit:      params.Limit,
	}, nil
}
