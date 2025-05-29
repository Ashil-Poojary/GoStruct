package models

import "time"

type RefreshToken struct {
	ID               int64     `db:"id"`
	UserID           int64     `db:"user_id"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	UserAgent        string    `db:"user_agent"`
	IPAddress        string    `db:"ip_address"`
	IssuedAt         time.Time `db:"issued_at"`
	ExpiresAt        time.Time `db:"expires_at"`
	Revoked          bool      `db:"revoked"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
