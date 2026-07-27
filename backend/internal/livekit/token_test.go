package livekit

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_ValidJWT(t *testing.T) {
	secret := "test-secret"
	tokenStr, err := GenerateToken("test-key", secret, "room-1", "user-1", "Alice", "recruiter", "interview-1")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("token is empty")
	}

	// Token must be a real, signature-valid JWT — not a mock string.
	parsed, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Fatalf("unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse/verify token: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("token is not valid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims are not MapClaims")
	}
	if claims["iss"] != "test-key" {
		t.Errorf("iss = %v, want test-key", claims["iss"])
	}
	if claims["sub"] != "user-1" {
		t.Errorf("sub = %v, want user-1", claims["sub"])
	}
	if claims["room_id"] != "room-1" {
		t.Errorf("room_id = %v, want room-1", claims["room_id"])
	}
	video, ok := claims["video"].(map[string]interface{})
	if !ok {
		t.Fatal("video grant missing")
	}
	if video["room"] != "room-1" || video["roomJoin"] != true {
		t.Errorf("video grant incorrect: %+v", video)
	}
}

func TestGenerateToken_FallbackSecrets(t *testing.T) {
	// Empty key/secret should fall back to dev defaults, still producing a
	// verifiable token (signed with devsecret).
	tokenStr, err := GenerateToken("", "", "room-x", "user-x", "Bob", "participant", "interview-x")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	_, err = jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		return []byte("devsecret"), nil
	})
	if err != nil {
		t.Fatalf("token not verifiable with devsecret fallback: %v", err)
	}
}
