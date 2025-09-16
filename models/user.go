package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug" gorm:"uniqueIndex"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Image    string `json:"image"`
	About    string `json:"about"`
}
