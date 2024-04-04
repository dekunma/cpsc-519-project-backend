package models

type Friendship struct {
	UserID   uint `json:"user_id" gorm:"primary_key; autoIncrement:false"`
	FriendID uint `json:"friend_id" gorm:"primary_key; autoIncrement:false"`
	Accepted bool `json:"accepted"`
}

type CreateFriendInvitationRequest struct {
	UserEmail   string `json:"user_email"`
	FriendEmail string `json:"friend_email"`
}
