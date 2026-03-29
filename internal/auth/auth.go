package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/huyng/anchat/internal/crypto"
	"github.com/huyng/anchat/internal/db"
	"github.com/huyng/anchat/internal/models"
)

// AuthService handles authentication and session management
type AuthService struct {
	db *db.DB
}

// NewAuthService creates a new authentication service
func NewAuthService(database *db.DB) *AuthService {
	return &AuthService{db: database}
}

// Session represents an active user session
type Session struct {
	UserID    string
	Token      string
	ExpiresAt time.Time
}

// RegisterUser registers a new user account
func (a *AuthService) RegisterUser(ctx context.Context, username, password string, pubkeyEd25519, pubkeyX25519 []byte) (*models.User, error) {
	// Generate salt and hash password
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	passwordHash, err := crypto.HashPassword(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Store salt with password hash (prepend salt to hash)
	passwordHashWithSalt := append(salt, passwordHash...)

	// Check if username already exists
	usernameHash := models.HashUsername(username)
	existing, err := a.db.GetUserByUsername(ctx, usernameHash)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("username already taken")
	}

	// Create user
	user := &models.User{
		UserID:        generateUserID(),
		UsernameHash:   usernameHash,
		PubkeyEd25519:  pubkeyEd25519,
		PubkeyX25519:  pubkeyX25519,
		PasswordHash:   passwordHashWithSalt,
		CreatedAt:      time.Now(),
	}

	if err := a.db.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// AuthenticateUser authenticates a user with username and password
func (a *AuthService) AuthenticateUser(ctx context.Context, username, password string, pubkeyEd25519, pubkeyX25519 []byte) (*Session, error) {
	usernameHash := models.HashUsername(username)
	user, err := a.db.GetUserByUsername(ctx, usernameHash)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Extract salt from stored hash (first 16 bytes)
	salt := user.PasswordHash[:16]
	storedHash := user.PasswordHash[16:]

	// Verify password
	valid, err := crypto.VerifyPassword(password, storedHash, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate session token
	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Store session token hash
	sessionTokenHash := models.HashSessionToken(sessionToken)
	if err := a.db.UpdateSessionToken(ctx, user.UserID, sessionTokenHash); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	session := &Session{
		UserID:    user.UserID,
		Token:      sessionToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
	}

	return session, nil
}

// AuthenticateBot authenticates a bot with API token
func (a *AuthService) AuthenticateBot(ctx context.Context, token string, pubkeyEd25519, pubkeyX25519 []byte) (*Session, error) {
	// TODO: Implement bot authentication
	return nil, fmt.Errorf("bot authentication not implemented")
}

// ValidateSession validates a session token
func (a *AuthService) ValidateSession(ctx context.Context, sessionToken string) (*models.User, error) {
	sessionTokenHash := models.HashSessionToken(sessionToken)
	user, err := a.db.GetUserBySessionToken(ctx, sessionTokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}
	return user, nil
}

// Logout invalidates a session
func (a *AuthService) Logout(ctx context.Context, userID string) error {
	return a.db.UpdateSessionToken(ctx, userID, nil)
}

// generateUserID generates a unique user ID
func generateUserID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// generateSessionToken generates a session token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
