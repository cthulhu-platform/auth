package service

import (
	"fmt"

	"github.com/cthulhu-platform/auth/internal/pkg"
)

type authenticationService struct{}

// NewAuthenticationService creates a new authentication service instance
func NewAuthenticationService() AuthenticationService {
	return &authenticationService{}
}

func (s *authenticationService) GenerateTokens(userID, email, provider string) (*pkg.TokenPair, error) {
	return nil, fmt.Errorf("token generation not implemented")
}

func (s *authenticationService) ValidateAccessToken(token string) (*pkg.Claims, error) {
	return nil, fmt.Errorf("token validation not implemented")
}

func (s *authenticationService) ValidateUserID(userID string) (bool, error) {
	return false, fmt.Errorf("user ID validation not implemented")
}

func (s *authenticationService) HandleDiagnoseMessage(transactionID, operation, message string) error {
	// Handle the diagnose message
	// The logging is done in the handler layer
	return nil
}

