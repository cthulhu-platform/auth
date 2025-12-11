package service

import "github.com/cthulhu-platform/auth/internal/pkg"

// AuthenticationService defines the interface for authentication operations
type AuthenticationService interface {
	GenerateTokens(userID, email, provider string) (*pkg.TokenPair, error)
	ValidateAccessToken(token string) (*pkg.Claims, error)
	ValidateUserID(userID string) (bool, error)
	HandleDiagnoseMessage(transactionID, operation, message string) error
}

