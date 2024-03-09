package entities

import "time"

type User struct {
	Id           uint32    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	IsDeleted    bool      `db:"is_deleted"`
	IsDisabled   bool      `db:"is_disabled"`
	PasswordHash string    `db:"password_hash"`
}

type UserSession struct {
	Id           uint32    `db:"id"`
	Uid          uint32    `db:"uid"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
	Valid        bool      `db:"valid"`
}
