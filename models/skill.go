package models

import "time"

type Skill struct {
	SkillID           uint64    `gorm:"primaryKey;autoIncrement" json:"skill_id"`
	CvID              uint64    `json:"cv_id"`
	SkillCategoryName string    `json:"skill_category_name"`
	SkillName         string    `json:"skill_name"`
	Desc              string    `json:"desc"`
	Icon              string    `json:"icon"`
	Nilai             int       `json:"nilai"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
