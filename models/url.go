package models

import "time"

type Url struct {
	Id         uint32     `json:"id"`
	Uid        uint32     `json:"uid"`
	Name       *string    `json:"name,omitempty"`
	ValidUntil *time.Time `json:"validUntil,omitempty"`
	Url        string     `json:"url"`
	RedirectTo string     `json:"redirectTo"`
}
