 // Package copilot provides authentication and token management for GitHub Copilot API.
 package copilot
 
 import (
 	"context"
 	"encoding/json"
 	"fmt"
 	"io"
 	"net/http"
 	"net/url"
 	"strings"
 	"time"
 
 	"github.com/radityprtama/proxygate/v6/internal/config"
 	"github.com/radityprtama/proxygate/v6/internal/util"
 	log "github.com/sirupsen/logrus"
 )
 
 const (
 	// GitHub OAuth endpoints
 	githubDeviceCodeURL = "https://github.com/login/device/code"
 	githubTokenURL      = "https://github.com/login/oauth/access_token"
 	githubUserURL       = "https://api.github.com/user"
 
 	// Copilot OAuth client ID (VS Code's client ID for Copilot)
 	copilotClientID = "Iv1.b507a08c87ecfe98"
 
 	// Required scopes for Copilot access
 	copilotScopes = "read:user"
 )
 
 // DeviceFlowClient handles GitHub OAuth device flow.
 type DeviceFlowClient struct {
 	httpClient *http.Client
 	cfg        *config.Config
 }
 
 // NewDeviceFlowClient creates a new device flow client.
 func NewDeviceFlowClient(cfg *config.Config) *DeviceFlowClient {
 	return &DeviceFlowClient{
 		httpClient: util.SetProxy(&cfg.SDKConfig, &http.Client{Timeout: 30 * time.Second}),
 		cfg:        cfg,
 	}
 }
 
 // RequestDeviceCode initiates the device flow and returns the device code response.
 func (c *DeviceFlowClient) RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, error) {
 	data := url.Values{}
 	data.Set("client_id", copilotClientID)
 	data.Set("scope", copilotScopes)
 
 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubDeviceCodeURL, strings.NewReader(data.Encode()))
 	if err != nil {
 		return nil, NewAuthenticationError(ErrDeviceFlowFailed, err)
 	}
 
 	req.Header.Set("Accept", "application/json")
 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
 
 	resp, err := c.httpClient.Do(req)
 	if err != nil {
 		return nil, NewAuthenticationError(ErrDeviceFlowFailed, err)
 	}
 	defer func() {
 		if errClose := resp.Body.Close(); errClose != nil {
 			log.Errorf("device code: close body error: %v", errClose)
 		}
 	}()
 
 	if resp.StatusCode != http.StatusOK {
 		body, _ := io.ReadAll(resp.Body)
 		return nil, NewAuthenticationError(ErrDeviceFlowFailed,
 			fmt.Errorf("status %d: %s", resp.StatusCode, string(body)))
 	}
 
 	var deviceCode DeviceCodeResponse
 	if err = json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
 		return nil, NewAuthenticationError(ErrDeviceFlowFailed, err)
 	}
 
 	if deviceCode.DeviceCode == "" {
 		return nil, NewAuthenticationError(ErrDeviceFlowFailed, fmt.Errorf("empty device code"))
 	}
 
 	return &deviceCode, nil
 }
 
 // PollForToken polls the token endpoint until the user authorizes or times out.
 func (c *DeviceFlowClient) PollForToken(ctx context.Context, deviceCode *DeviceCodeResponse) (*TokenResponse, error) {
 	interval := time.Duration(deviceCode.Interval) * time.Second
 	if interval < 5*time.Second {
 		interval = 5 * time.Second
 	}
 
 	deadline := time.Now().Add(time.Duration(deviceCode.ExpiresIn) * time.Second)
 
 	for time.Now().Before(deadline) {
 		select {
 		case <-ctx.Done():
 			return nil, ctx.Err()
 		case <-time.After(interval):
 		}
 
 		token, err := c.exchangeDeviceCode(ctx, deviceCode.DeviceCode)
 		if err != nil {
 			var authErr *AuthenticationError
 			if ok := isAuthError(err, &authErr); ok {
 				switch {
 				case authErr.Type == ErrAuthorizationPending:
 					continue
 				case authErr.Type == ErrSlowDown:
 					interval += 5 * time.Second
 					continue
 				}
 			}
 			return nil, err
 		}
 
 		return token, nil
 	}
 
 	return nil, NewAuthenticationError(ErrExpiredToken, fmt.Errorf("device code expired"))
 }
 
 func isAuthError(err error, target **AuthenticationError) bool {
 	if authErr, ok := err.(*AuthenticationError); ok {
 		*target = authErr
 		return true
 	}
 	return false
 }
 
 // exchangeDeviceCode exchanges the device code for an access token.
 func (c *DeviceFlowClient) exchangeDeviceCode(ctx context.Context, deviceCode string) (*TokenResponse, error) {
 	data := url.Values{}
 	data.Set("client_id", copilotClientID)
 	data.Set("device_code", deviceCode)
 	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
 
 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubTokenURL, strings.NewReader(data.Encode()))
 	if err != nil {
 		return nil, NewAuthenticationError(ErrTokenExchangeFailed, err)
 	}
 
 	req.Header.Set("Accept", "application/json")
 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
 
 	resp, err := c.httpClient.Do(req)
 	if err != nil {
 		return nil, NewAuthenticationError(ErrTokenExchangeFailed, err)
 	}
 	defer func() {
 		if errClose := resp.Body.Close(); errClose != nil {
 			log.Errorf("token exchange: close body error: %v", errClose)
 		}
 	}()
 
 	var token TokenResponse
 	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
 		return nil, NewAuthenticationError(ErrTokenExchangeFailed, err)
 	}
 
 	if token.Error != "" {
 		switch token.Error {
 		case "authorization_pending":
 			return nil, &AuthenticationError{Type: ErrAuthorizationPending}
 		case "slow_down":
 			return nil, &AuthenticationError{Type: ErrSlowDown}
 		case "expired_token":
 			return nil, &AuthenticationError{Type: ErrExpiredToken}
 		case "access_denied":
 			return nil, &AuthenticationError{Type: ErrAccessDenied, Message: token.ErrorDescription}
 		default:
 			return nil, NewAuthenticationError(ErrTokenExchangeFailed,
 				fmt.Errorf("%s: %s", token.Error, token.ErrorDescription))
 		}
 	}
 
 	if token.AccessToken == "" {
 		return nil, NewAuthenticationError(ErrTokenExchangeFailed, fmt.Errorf("empty access token"))
 	}
 
 	return &token, nil
 }
 
 // FetchUserInfo fetches the GitHub user information using the access token.
 func (c *DeviceFlowClient) FetchUserInfo(ctx context.Context, accessToken string) (string, error) {
 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
 	if err != nil {
 		return "", err
 	}
 
 	req.Header.Set("Authorization", "token "+accessToken)
 	req.Header.Set("Accept", "application/json")
 
 	resp, err := c.httpClient.Do(req)
 	if err != nil {
 		return "", err
 	}
 	defer func() {
 		if errClose := resp.Body.Close(); errClose != nil {
 			log.Errorf("fetch user info: close body error: %v", errClose)
 		}
 	}()
 
 	if resp.StatusCode != http.StatusOK {
 		return "", fmt.Errorf("failed to fetch user info: status %d", resp.StatusCode)
 	}
 
 	var user struct {
 		Login string `json:"login"`
 	}
 	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
 		return "", err
 	}
 
 	return user.Login, nil
 }
