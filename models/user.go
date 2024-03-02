package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
}

type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type CheckVerificationCodeRequest struct {
	Email            string `json:"email" binding:"required"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

type SignUpRequest struct {
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verification_code"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
