package model

import "time"

type Category struct {
	ID        string
	UserID    string
	Name      string
	NameKey   *string
	Color     string
	SortOrder int
	IsVisible bool
	IsSystem  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
