 // Package copilot provides authentication and token management for GitHub Copilot API.
 package copilot
 
 // CopilotTokenStorage represents the stored token data for GitHub Copilot.
 type CopilotTokenStorage struct {
 	// AccessToken is the GitHub OAuth access token.
 	AccessToken string `json:"access_token"`
 	// TokenType is the type of token (usually "bearer").
 	TokenType string `json:"token_type"`
 	// Scope is the OAuth scope granted.
 	Scope string `json:"scope"`
 	// Username is the GitHub username.
 	Username string `json:"username,omitempty"`
 	// Type identifies this as a GitHub Copilot token.
 	Type string `json:"type"`
 }
 
 // DeviceCodeResponse represents the response from GitHub's device code endpoint.
 type DeviceCodeResponse struct {
 	// DeviceCode is the device verification code.
 	DeviceCode string `json:"device_code"`
 	// UserCode is the code the user enters at the verification URL.
 	UserCode string `json:"user_code"`
 	// VerificationURI is the URL where the user enters the code.
 	VerificationURI string `json:"verification_uri"`
 	// ExpiresIn is the lifetime in seconds of the device code.
 	ExpiresIn int `json:"expires_in"`
 	// Interval is the minimum seconds between polling requests.
 	Interval int `json:"interval"`
 }
 
 // TokenResponse represents the response from GitHub's token endpoint.
 type TokenResponse struct {
 	// AccessToken is the OAuth access token.
 	AccessToken string `json:"access_token"`
 	// TokenType is the type of token.
 	TokenType string `json:"token_type"`
 	// Scope is the OAuth scope granted.
 	Scope string `json:"scope"`
 	// Error is set if the request failed.
 	Error string `json:"error,omitempty"`
 	// ErrorDescription provides more details about the error.
 	ErrorDescription string `json:"error_description,omitempty"`
 }
 
 // CopilotAuthBundle contains the authentication data after successful auth.
 type CopilotAuthBundle struct {
 	// TokenData contains the OAuth token information.
 	TokenData *TokenResponse
 	// Username is the GitHub username.
 	Username string
 }
