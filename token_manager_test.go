package gigachat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewTokenManager(t *testing.T) {
	authKey := "test_auth_key"
	tm := NewTokenManager(authKey)
	
	if tm.authKey != authKey {
		t.Errorf("Expected authKey %v, got %v", authKey, tm.authKey)
	}
	if tm.scope != "GIGACHAT_API_PERS" {
		t.Errorf("Expected default scope GIGACHAT_API_PERS, got %v", tm.scope)
	}
	if tm.oauthURI != "https://ngw.devices.sberbank.ru:9443" {
		t.Errorf("Expected default oauthURI, got %v", tm.oauthURI)
	}
	if tm.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestNewTokenManager_WithOptions(t *testing.T) {
	authKey := "test_auth_key"
	customScope := "CUSTOM_SCOPE"
	customURI := "https://custom.uri"
	
	tm := NewTokenManager(
		authKey,
		WithScope(customScope),
		WithOAuthURI(customURI),
		WithInsecureSkipVerify(true),
	)
	
	if tm.scope != customScope {
		t.Errorf("Expected scope %v, got %v", customScope, tm.scope)
	}
	if tm.oauthURI != customURI {
		t.Errorf("Expected oauthURI %v, got %v", customURI, tm.oauthURI)
	}
	
	if transport, ok := tm.client.Transport.(*http.Transport); ok {
		if !transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be true")
		}
	} else {
		t.Error("Expected HTTP transport to be *http.Transport")
	}
}

func TestTokenManager_GetAccessToken_Cached(t *testing.T) {
	tm := NewTokenManager("test_key")
	
	tm.accessToken = "cached_token"
	tm.expiresAt = time.Now().Add(5 * time.Minute)
	
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token != "cached_token" {
		t.Errorf("Expected cached token, got %v", token)
	}
}

func TestTokenManager_GetAccessToken_Refresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/oauth" {
			t.Errorf("Expected path /api/v2/oauth, got %v", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %v", r.Method)
		}
		
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			t.Errorf("Expected Basic auth header, got %v", authHeader)
		}
		
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Expected application/x-www-form-urlencoded, got %v", contentType)
		}
		
		response := TokenResponse{
			AccessToken: "new_token_123",
			ExpiresAt:   time.Now().Add(30 * time.Minute).Unix(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_auth_key",
		WithOAuthURI(server.URL),
	)
	
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token != "new_token_123" {
		t.Errorf("Expected token new_token_123, got %v", token)
	}
	if tm.accessToken != "new_token_123" {
		t.Errorf("Expected cached token new_token_123, got %v", tm.accessToken)
	}
}

func TestTokenManager_GetAccessToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid credentials"}`))
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"invalid_key",
		WithOAuthURI(server.URL),
	)
	
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("Expected AuthenticationError, got %T", err)
	}
}

func TestTokenManager_GetAccessToken_EmptyToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TokenResponse{
			AccessToken: "",
			ExpiresAt:   time.Now().Add(30 * time.Minute).Unix(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_key",
		WithOAuthURI(server.URL),
	)
	
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error for empty token, got nil")
	}
	
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("Expected AuthenticationError, got %T", err)
	}
}

func TestTokenManager_GetAccessToken_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_key",
		WithOAuthURI(server.URL),
	)
	
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
	
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("Expected AuthenticationError, got %T", err)
	}
}

func TestTokenManager_ExpiresAt_Milliseconds(t *testing.T) {
	expiresAtMs := time.Now().Add(30 * time.Minute).UnixMilli()
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TokenResponse{
			AccessToken: "test_token",
			ExpiresAt:   expiresAtMs,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_key",
		WithOAuthURI(server.URL),
	)
	
	_, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expectedExpiresAt := time.Unix(expiresAtMs/1000, 0)
	if !tm.expiresAt.Equal(expectedExpiresAt) {
		t.Errorf("Expected expiresAt %v, got %v", expectedExpiresAt, tm.expiresAt)
	}
}

func TestTokenManager_ExpiresAt_NoExpiresAt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TokenResponse{
			AccessToken: "test_token",
			ExpiresAt:   0,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_key",
		WithOAuthURI(server.URL),
	)
	
	before := time.Now()
	_, err := tm.GetAccessToken()
	after := time.Now()
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expectedMin := before.Add(29 * time.Minute)
	expectedMax := after.Add(29 * time.Minute)
	
	if tm.expiresAt.Before(expectedMin) || tm.expiresAt.After(expectedMax) {
		t.Errorf("Expected expiresAt between %v and %v, got %v", expectedMin, expectedMax, tm.expiresAt)
	}
}

func TestTokenManager_ConcurrentAccess(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		time.Sleep(50 * time.Millisecond)
		
		response := TokenResponse{
			AccessToken: "concurrent_token",
			ExpiresAt:   time.Now().Add(30 * time.Minute).Unix(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := NewTokenManager(
		"test_key",
		WithOAuthURI(server.URL),
	)
	
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := tm.GetAccessToken()
			if err != nil {
				t.Errorf("Concurrent access error: %v", err)
			}
			done <- true
		}()
	}
	
	for i := 0; i < 5; i++ {
		<-done
	}
	
	if callCount > 2 {
		t.Errorf("Expected at most 2 API calls due to locking, got %d", callCount)
	}
}

func TestWithTokenManagerHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	tm := NewTokenManager(
		"test_key",
		WithTokenManagerHTTPClient(customClient),
	)
	
	if tm.client != customClient {
		t.Error("Expected custom HTTP client to be set")
	}
}
