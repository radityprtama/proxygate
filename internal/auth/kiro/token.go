package kiro

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// KiroTokenData holds OAuth token information from AWS CodeWhisperer (Kiro)
type KiroTokenData struct {
	// AccessToken is the OAuth2 access token for API access
	AccessToken string `json:"accessToken"`
	// RefreshToken is used to obtain new access tokens
	RefreshToken string `json:"refreshToken"`
	// ProfileArn is the AWS CodeWhisperer profile ARN
	ProfileArn string `json:"profileArn"`
	// ExpiresAt is the timestamp when the token expires
	ExpiresAt string `json:"expiresAt"`
	// AuthMethod indicates the authentication method used (e.g., "builder-id", "social")
	AuthMethod string `json:"authMethod"`
	// Provider indicates the OAuth provider (e.g., "AWS", "Google")
	Provider string `json:"provider"`
	// ClientID is the OIDC client ID (needed for token refresh)
	ClientID string `json:"clientId,omitempty"`
	// ClientSecret is the OIDC client secret (needed for token refresh)
	ClientSecret string `json:"clientSecret,omitempty"`
	// Email is the user's email address (used for file naming)
	Email string `json:"email,omitempty"`
	// StartURL is the IDC/Identity Center start URL (only for IDC auth method)
	StartURL string `json:"startUrl,omitempty"`
	// Region is the AWS region for IDC authentication (only for IDC auth method)
	Region string `json:"region,omitempty"`
}

// KiroTokenStorage holds the persistent token data for Kiro authentication.
type KiroTokenStorage struct {
	// AccessToken is the OAuth2 access token for API access
	AccessToken string `json:"access_token"`
	// RefreshToken is used to obtain new access tokens
	RefreshToken string `json:"refresh_token"`
	// ProfileArn is the AWS CodeWhisperer profile ARN
	ProfileArn string `json:"profile_arn"`
	// ExpiresAt is the timestamp when the token expires
	ExpiresAt string `json:"expires_at"`
	// AuthMethod indicates the authentication method used
	AuthMethod string `json:"auth_method"`
	// Provider indicates the OAuth provider
	Provider string `json:"provider"`
	// ClientID is the OIDC client ID
	ClientID string `json:"client_id,omitempty"`
	// ClientSecret is the OIDC client secret
	ClientSecret string `json:"client_secret,omitempty"`
	// Email is the user's email address
	Email string `json:"email,omitempty"`
	// StartURL is the IDC start URL
	StartURL string `json:"start_url,omitempty"`
	// Region is the AWS region for IDC
	Region string `json:"region,omitempty"`
}

// JWTClaims represents the claims we care about from a JWT token.
type JWTClaims struct {
	Email          string `json:"email,omitempty"`
	Sub            string `json:"sub,omitempty"`
	PreferredUser  string `json:"preferred_username,omitempty"`
	Name           string `json:"name,omitempty"`
	Iss            string `json:"iss,omitempty"`
}

// KiroIDETokenFile is the default path to Kiro IDE's token file
const KiroIDETokenFile = ".aws/sso/cache/kiro-auth-token.json"

// LoadKiroIDEToken loads token data from Kiro IDE's token file.
func LoadKiroIDEToken() (*KiroTokenData, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	tokenPath := filepath.Join(homeDir, KiroIDETokenFile)
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Kiro IDE token file (%s): %w", tokenPath, err)
	}

	var token KiroTokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse Kiro IDE token: %w", err)
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("access token is empty in Kiro IDE token file")
	}

	return &token, nil
}

// ExtractEmailFromJWT extracts the user's email from a JWT access token.
// JWT tokens typically have format: header.payload.signature
func ExtractEmailFromJWT(accessToken string) string {
	if accessToken == "" {
		return ""
	}

	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return ""
	}

	payload := parts[1]

	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return ""
		}
	}

	var claims JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return ""
	}

	if claims.Email != "" {
		return claims.Email
	}

	if claims.PreferredUser != "" && strings.Contains(claims.PreferredUser, "@") {
		return claims.PreferredUser
	}

	if claims.Sub != "" && strings.Contains(claims.Sub, "@") {
		return claims.Sub
	}

	return ""
}

// SanitizeEmailForFilename sanitizes an email address for use in a filename.
func SanitizeEmailForFilename(email string) string {
	if email == "" {
		return ""
	}

	result := email

	result = strings.ReplaceAll(result, "%2F", "_")
	result = strings.ReplaceAll(result, "%2f", "_")
	result = strings.ReplaceAll(result, "%5C", "_")
	result = strings.ReplaceAll(result, "%5c", "_")
	result = strings.ReplaceAll(result, "%2E", "_")
	result = strings.ReplaceAll(result, "%2e", "_")
	result = strings.ReplaceAll(result, "%00", "_")
	result = strings.ReplaceAll(result, "%", "_")

	for _, char := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " ", "\x00"} {
		result = strings.ReplaceAll(result, char, "_")
	}

	parts := strings.Split(result, "_")
	for i, part := range parts {
		for strings.HasPrefix(part, ".") {
			part = "_" + part[1:]
		}
		parts[i] = part
	}
	result = strings.Join(parts, "_")

	return result
}

// ToTokenData converts storage to KiroTokenData for API use.
func (s *KiroTokenStorage) ToTokenData() *KiroTokenData {
	return &KiroTokenData{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		ProfileArn:   s.ProfileArn,
		ExpiresAt:    s.ExpiresAt,
		AuthMethod:   s.AuthMethod,
		Provider:     s.Provider,
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		Email:        s.Email,
		StartURL:     s.StartURL,
		Region:       s.Region,
	}
}
