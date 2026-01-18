package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TokenManager struct {
	authKey     string
	scope       string
	oauthURI    string
	client      *http.Client
	accessToken string
	expiresAt   time.Time
	mu          sync.RWMutex
}

func NewTokenManager(authKey string, options ...TokenManagerOption) *TokenManager {
	tm := &TokenManager{
		authKey:  authKey,
		scope:    "GIGACHAT_API_PERS",
		oauthURI: "https://ngw.devices.sberbank.ru:9443",
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
	}

	for _, opt := range options {
		opt(tm)
	}

	return tm
}

type TokenManagerOption func(*TokenManager)

func WithScope(scope string) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.scope = scope
	}
}

func WithOAuthURI(uri string) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.oauthURI = uri
	}
}

func WithInsecureSkipVerify(skip bool) TokenManagerOption {
	return func(tm *TokenManager) {
		if transport, ok := tm.client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

func WithTokenManagerHTTPClient(client *http.Client) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.client = client
	}
}

func (tm *TokenManager) GetAccessToken() (string, error) {
	tm.mu.RLock()
	if tm.accessToken != "" && time.Now().Add(30*time.Second).Before(tm.expiresAt) {
		token := tm.accessToken
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	return tm.refreshToken()
}

func (tm *TokenManager) refreshToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.accessToken != "" && time.Now().Add(30*time.Second).Before(tm.expiresAt) {
		return tm.accessToken, nil
	}

	rqUID := uuid.New().String()

	formData := url.Values{}
	formData.Set("scope", tm.scope)

	req, err := http.NewRequest("POST", tm.oauthURI+"/api/v2/oauth", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", &AuthenticationError{Message: "failed to create request", Err: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", rqUID)
	req.Header.Set("Authorization", "Basic "+tm.authKey)

	resp, err := tm.client.Do(req)
	if err != nil {
		return "", &AuthenticationError{Message: "OAuth token request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", &AuthenticationError{
			Message: fmt.Sprintf("OAuth request failed with status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", &AuthenticationError{Message: "failed to decode token response", Err: err}
	}

	if tokenResp.AccessToken == "" {
		return "", &AuthenticationError{Message: "invalid token response: empty access token"}
	}

	tm.accessToken = tokenResp.AccessToken

	if tokenResp.ExpiresAt > 0 {
		expiresAt := tokenResp.ExpiresAt
		if expiresAt > 1000000000000 {
			expiresAt = expiresAt / 1000
		}
		tm.expiresAt = time.Unix(expiresAt, 0)
	} else {
		tm.expiresAt = time.Now().Add(29 * time.Minute)
	}

	return tm.accessToken, nil
}
