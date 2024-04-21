package entity

import "time"

type Session struct {
	GUID             string    `json:"guid" bson:"guid"`
	RefreshTokenHash string    `bson:"refresh_token_hash"`
	ExpiresAt        time.Time `bson:"expires_at"`
}
