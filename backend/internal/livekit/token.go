package livekit

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a LiveKit compatible JWT token for room access.
func GenerateToken(apiKey, apiSecret, roomName, identity, displayName, role, interviewID string) (string, error) {
	if apiSecret == "" {
		apiSecret = os.Getenv("LIVEKIT_API_SECRET")
		if apiSecret == "" {
			apiSecret = "devsecret"
		}
	}
	if apiKey == "" {
		apiKey = os.Getenv("LIVEKIT_API_KEY")
		if apiKey == "" {
			apiKey = "devkey"
		}
	}

	claims := jwt.MapClaims{
		"iss":          apiKey,
		"sub":          identity,
		"nbf":          time.Now().Unix() - 5,
		"exp":          time.Now().Add(4 * time.Hour).Unix(), // valid for 4 hours
		"name":         displayName,
		"role":         role,
		"room_id":      roomName,
		"interview_id": interviewID,
		"video": map[string]interface{}{
			"room":         roomName,
			"roomJoin":     true,
			"canPublish":   true,
			"canSubscribe": true,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
