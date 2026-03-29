package models

import (
	"crypto/sha256"
	"strings"
	"time"
)

// Channel represents a chat channel
type Channel struct {
	ChannelID   string    `json:"channel_id" db:"channel_id"`
	NameHash    []byte    `json:"-" db:"name_hash"`
	MemberCount int       `json:"member_count" db:"member_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ChannelMember represents a user in a channel
type ChannelMember struct {
	ChannelID string    `json:"channel_id" db:"channel_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
	IsOp      bool      `json:"is_op" db:"is_op"` // channel operator
}

// HashChannelName normalizes and hashes a channel name
func HashChannelName(name string) []byte {
	// Normalize: lowercase, trim whitespace
	normalized := strings.ToLower(strings.TrimSpace(name))
	hash := sha256.Sum256([]byte(normalized))
	return hash[:]
}
