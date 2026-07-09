package jwt

import (
	"testing"
	"time"
)

const testSecret = "test-secret-key-1234567890"

func TestGenerateTokenPair(t *testing.T) {
	access, refresh, err := GenerateTokenPair("user-1", "a@b.com", "recruiter", testSecret, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}
	if access == "" {
		t.Fatal("access token must not be empty")
	}
	if refresh == "" {
		t.Fatal("refresh token must not be empty")
	}
}

func TestValidateToken_Success(t *testing.T) {
	access, _, err := GenerateTokenPair("user-1", "a@b.com", "recruiter", testSecret, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateToken(access, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("expected user-1, got %s", claims.UserID)
	}
	if claims.Email != "a@b.com" {
		t.Fatalf("expected a@b.com, got %s", claims.Email)
	}
	if claims.Role != "recruiter" {
		t.Fatalf("expected recruiter, got %s", claims.Role)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	expired, _, _ := GenerateTokenPair("user-1", "a@b.com", "recruiter", testSecret, -1*time.Hour, 7*24*time.Hour)
	_, err := ValidateToken(expired, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	access, _, _ := GenerateTokenPair("user-1", "a@b.com", "recruiter", "other-secret", 15*time.Minute, 7*24*time.Hour)
	_, err := ValidateToken(access, testSecret)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	_, err := ValidateToken("invalid-token-string", testSecret)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateToken_RefreshToken(t *testing.T) {
	_, refresh, err := GenerateTokenPair("user-1", "a@b.com", "recruiter", testSecret, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateToken(refresh, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken on refresh token: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("expected user-1, got %s", claims.UserID)
	}
}
