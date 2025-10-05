package models

type Post struct {
	PostID     int64  `gorm:"column:post_id;primaryKey" json:"post_id"`
	UserID     int64  `gorm:"column:user_id" json:"user_id"`
	CategoryID int64  `gorm:"column:category_id" json:"category_id"`
	Title      string `gorm:"column:title" json:"title"`
	Slug       string `gorm:"column:slug" json:"slug"`
	Thumbnail  string `gorm:"column:thumbnail" json:"thumbnail"`
	Excerpt    string `gorm:"column:excerpt" json:"excerpt"`
	Content    string `gorm:"column:content" json:"content"`
	Status     int    `gorm:"column:status" json:"status"`
	Views      int    `gorm:"column:views" json:"views"`
}

func (Post) TableName() string {
	return ("posts")
}

type PostCatgeory struct {
	CategoryID int64  `gorm:"column:category_id;primaryKey" json:"category_id"`
	Name       string `gorm:"column:name" json:"name"`
	Slug       string `gorm:"column:slug" json:"slug"`
	Desc       string `gorm:"column:description" json:"description"`
}

func (PostCatgeory) TableName() string {
	return ("post_categories")
}

type Tags struct {
	CategoryID int64  `gorm:"column:tag_id;primaryKey" json:"tag_id"`
	Name       string `gorm:"column:name" json:"name"`
	Slug       string `gorm:"column:slug" json:"slug"`
}

func (Tags) TableName() string {
	return ("tags")
}

type PostTags struct {
	ID     int64 `gorm:"column:id;primaryKey" json:"id"`
	TagID  int64 `gorm:"column:tag_id" json:"tag_id"`
	PostID int64 `gorm:"column:post_id" json:"post_id"`
}

func (PostTags) TableName() string {
	return ("post_tags")
}
