package model

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Gender    Gender     `json:"gender"`
	BirthDate time.Time  `json:"birth_date"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type CursorPaginatedResponse struct {
	Data       []User `json:"data"`
	NextCursor *int   `json:"next_cursor,omitempty"`
	Limit      int    `json:"limit"`
}
