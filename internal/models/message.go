package models

import "time"

// Message represents a stored message (encrypted at rest)
type Message struct {
	ID             int64     `json:"id" db:"id"`
	ChannelID      *string   `json:"-" db:"channel_id"`      // NULL for private messages
	RecipientID    *string   `json:"-" db:"recipient_user_id"` // NULL for channel messages
	SenderKeyHash  []byte    `json:"-" db:"sender_key_hash"`  // blind index (SHA-256 of sender's pubkey)
	EncryptedBlob  []byte    `json:"-" db:"encrypted_blob"`    // already E2E encrypted by client
	Signature      []byte    `json:"-" db:"signature"`         // sender signs message
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
}

// ChannelMessage is sent to clients (includes metadata for display)
type ChannelMessage struct {
	MessageID  int64  `json:"message_id"`
	Channel    string `json:"channel"`
	From       string `json:"from"`
	Ciphertext string `json:"ciphertext"` // base64url encoded
	Nonce      string `json:"nonce"`      // base64url encoded
	Timestamp  int64  `json:"timestamp"`
}

// PrivateMessage is a private message event
type PrivateMessage struct {
	MessageID  int64  `json:"message_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Ciphertext string `json:"ciphertext"` // base64url encoded
	Nonce      string `json:"nonce"`      // base64url encoded
	Timestamp  int64  `json:"timestamp"`
}
