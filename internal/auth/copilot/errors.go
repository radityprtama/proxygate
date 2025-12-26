 // Package copilot provides authentication and token management for GitHub Copilot API.
 package copilot
 
 import (
 	"errors"
 	"fmt"
 )
 
 // Authentication error types
 var (
 	ErrDeviceFlowFailed    = errors.New("device flow failed")
 	ErrTokenExchangeFailed = errors.New("token exchange failed")
 	ErrAuthorizationPending = errors.New("authorization pending")
 	ErrSlowDown            = errors.New("slow down")
 	ErrAccessDenied        = errors.New("access denied")
 	ErrExpiredToken        = errors.New("expired token")
 	ErrUnsupportedGrant    = errors.New("unsupported grant type")
 )
 
 // AuthenticationError represents a Copilot authentication error.
 type AuthenticationError struct {
 	Type    error
 	Message string
 	Cause   error
 }
 
 // Error implements the error interface.
 func (e *AuthenticationError) Error() string {
 	if e.Cause != nil {
 		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
 	}
 	if e.Message != "" {
 		return fmt.Sprintf("%s: %s", e.Type, e.Message)
 	}
 	return e.Type.Error()
 }
 
 // Unwrap returns the underlying error.
 func (e *AuthenticationError) Unwrap() error {
 	return e.Type
 }
 
 // NewAuthenticationError creates a new AuthenticationError.
 func NewAuthenticationError(errType error, cause error) *AuthenticationError {
 	return &AuthenticationError{
 		Type:  errType,
 		Cause: cause,
 	}
 }
 
 // GetUserFriendlyMessage returns a user-friendly error message.
 func GetUserFriendlyMessage(err error) string {
 	var authErr *AuthenticationError
 	if errors.As(err, &authErr) {
 		switch {
 		case errors.Is(authErr.Type, ErrAccessDenied):
 			return "Access denied. Please ensure you have an active GitHub Copilot subscription."
 		case errors.Is(authErr.Type, ErrExpiredToken):
 			return "Authorization expired. Please try again."
 		case errors.Is(authErr.Type, ErrDeviceFlowFailed):
 			return "Failed to start device authorization. Please check your network connection."
 		case errors.Is(authErr.Type, ErrTokenExchangeFailed):
 			return "Failed to exchange token. Please try again."
 		}
 	}
 	return err.Error()
 }
