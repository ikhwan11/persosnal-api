package models

import "time"

type Education struct {
	EduId           uint64    `gorm:"primaryKey;autoIncrement" json:"edu_id"`
	CvID            uint64    `json:"cv_id"`
	UniversitasName string    `json:"universitas_name"`
	Jurusan         string    `json:"jurusan"`
	Tahun           string    `json:"tahun"`
	Desc            string    `json:"desc"`
	IPK             float64   `gorm:"column:ipk" json:"ipk"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
