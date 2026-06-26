package realtime

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims holds the data extracted from a valid room_access_token.
// The token is a LiveKit JWT that carries room-scoped claims.
type TokenClaims struct {
	UserID      string
	Role        string // "recruiter" | "candidate" | "ai" | "guest"
	DisplayName string
	RoomID      string
	InterviewID string
}

// extractAndValidateToken reads the room_access_token from the request
// (query param ?token=… or Authorization: Bearer header) and validates it.
func extractAndValidateToken(r *http.Request) (*TokenClaims, error) {
	raw := tokenFromRequest(r)
	if raw == "" {
		return nil, errors.New("missing room_access_token")
	}

	return validateRoomToken(raw)
}

// tokenFromRequest extracts the raw token string from:
//  1. Query parameter: ?token=<value>
//  2. Authorization header: Bearer <value>
func tokenFromRequest(r *http.Request) string {
	if t := r.URL.Query().Get("token"); t != "" {
		return t
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// validateRoomToken parses and validates the JWT.
// The LiveKit JWT signing key is read from the environment variable
// LIVEKIT_API_SECRET (set via Docker / .env).
func validateRoomToken(raw string) (*TokenClaims, error) {
	if raw == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		secret := os.Getenv("LIVEKIT_API_SECRET")
		if secret == "" {
			secret = "devsecret" // Default fallback for development
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Extract UserID from "user_id" or standard "sub" claim
	var userID string
	if uid, ok := claims["user_id"].(string); ok && uid != "" {
		userID = uid
	} else if sub, ok := claims["sub"].(string); ok && sub != "" {
		userID = sub
	}

	// Extract RoomID from "video.room" (LiveKit grant object) or custom "room_id" / "room" claim
	var roomID string
	if video, ok := claims["video"].(map[string]interface{}); ok {
		if r, ok := video["room"].(string); ok {
			roomID = r
		}
	}
	if roomID == "" {
		if r, ok := claims["room_id"].(string); ok && r != "" {
			roomID = r
		} else if r, ok := claims["room"].(string); ok && r != "" {
			roomID = r
		}
	}

	// Extract Role from "role"
	var role string
	if r, ok := claims["role"].(string); ok {
		role = r
	}

	// Extract DisplayName from "name" or standard "display_name" claim
	var displayName string
	if name, ok := claims["name"].(string); ok {
		displayName = name
	} else if name, ok := claims["display_name"].(string); ok {
		displayName = name
	}

	// Extract InterviewID from "interview_id"
	var interviewID string
	if iID, ok := claims["interview_id"].(string); ok {
		interviewID = iID
	}

	// Check required fields
	if userID == "" {
		return nil, errors.New("token is missing user identity (user_id/sub)")
	}
	if roomID == "" {
		return nil, errors.New("token is missing room identifier (room_id/video.room)")
	}
	if role == "" {
		return nil, errors.New("token is missing role claim")
	}

	// Normalize and validate role
	role = strings.ToLower(role)
	if role != "recruiter" && role != "candidate" && role != "ai" {
		return nil, fmt.Errorf("role '%s' is not authorized to join rooms", role)
	}

	return &TokenClaims{
		UserID:      userID,
		Role:        role,
		DisplayName: displayName,
		RoomID:      roomID,
		InterviewID: interviewID,
	}, nil
}

// extractAndValidateRecruiterToken validates a recruiter's Bearer token.
// It is used for API endpoints outside the WebSocket connection (e.g. requesting a room token),
// meaning the token does NOT need to contain a RoomID grant yet.
func extractAndValidateRecruiterToken(r *http.Request) (*TokenClaims, error) {
	raw := tokenFromRequest(r)
	if raw == "" {
		return nil, errors.New("missing token")
	}

	token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		secret := os.Getenv("LIVEKIT_API_SECRET")
		if secret == "" {
			secret = "devsecret"
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	var userID string
	if uid, ok := claims["user_id"].(string); ok && uid != "" {
		userID = uid
	} else if sub, ok := claims["sub"].(string); ok && sub != "" {
		userID = sub
	}

	var role string
	if rClaim, ok := claims["role"].(string); ok {
		role = rClaim
	}

	var displayName string
	if name, ok := claims["name"].(string); ok {
		displayName = name
	} else if name, ok := claims["display_name"].(string); ok {
		displayName = name
	}

	if userID == "" {
		return nil, errors.New("token is missing user identity (user_id/sub)")
	}
	if role == "" {
		return nil, errors.New("token is missing role claim")
	}

	role = strings.ToLower(role)
	if role != "recruiter" && role != "admin" {
		return nil, fmt.Errorf("role '%s' is not authorized to manage recruiter resources", role)
	}

	return &TokenClaims{
		UserID:      userID,
		Role:        role,
		DisplayName: displayName,
	}, nil
}

