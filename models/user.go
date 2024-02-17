package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type SignUpRequest struct {
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verification_code"`
}
