package models

import "time"

type CV struct {
	CvID      uint64    `gorm:"primaryKey;autoIncrement" json:"cv_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Image     string    `json:"image"`
	Tagline   string    `json:"tagline"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
