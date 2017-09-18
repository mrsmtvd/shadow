package auth

import (
	"time"

	"github.com/markbates/goth"
)

type User struct {
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
}

func NewUser(u goth.User) *User {
	return &User{
		Provider:          u.Provider,
		Email:             u.Email,
		Name:              u.Name,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		Description:       u.Description,
		UserID:            u.UserID,
		AvatarURL:         u.AvatarURL,
		Location:          u.Location,
		AccessToken:       u.AccessToken,
		AccessTokenSecret: u.AccessTokenSecret,
		RefreshToken:      u.RefreshToken,
		ExpiresAt:         u.ExpiresAt,
	}
}

func (u *User) IsAuthorized() bool {
	return u.UserID != ""
}
