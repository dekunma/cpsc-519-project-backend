package models

import (
	"time"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PostImages struct {
	ImageID   uint      `json:"image_id" gorm:"primary_key"`
	PostID    uint      `json:"post_id"`   // Foreign key - Post.ID
	FilePath  string    `json:"file_path"` // S3 file path
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type UploadPostImageRequest struct {
	PostID uint   `json:"post_id" binding:"required"`
	Image  string `json:"image" binding:"required"`
}
