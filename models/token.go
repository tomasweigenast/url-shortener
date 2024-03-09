package models

import "time"

type AccessToken struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Expires      time.Time `json:"expires"`
}
