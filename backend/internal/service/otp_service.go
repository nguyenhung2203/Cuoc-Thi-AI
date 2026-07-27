package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// OTPService issues and verifies short-lived one-time passwords, backed by Redis.
// Purposes namespace the keys so a registration OTP can't be replayed as a
// password-reset OTP and vice versa.
type OTPService struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewOTPService(rc *redis.Client) *OTPService {
	return &OTPService{redis: rc, ttl: 10 * time.Minute}
}

// Enabled reports whether a Redis backend is wired. Without it, OTP cannot be
// issued or verified and the caller should fail loudly rather than pretend.
func (s *OTPService) Enabled() bool {
	return s.redis != nil
}

func (s *OTPService) key(purpose, email string) string {
	return fmt.Sprintf("otp:%s:%s", purpose, email)
}

// Generate creates a 6-digit code, stores it under (purpose,email) with a TTL,
// and returns the raw code so the caller can email it.
func (s *OTPService) Generate(ctx context.Context, purpose, email string) (string, error) {
	if !s.Enabled() {
		return "", fmt.Errorf("otp backend (redis) not configured")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())
	if err := s.redis.Set(ctx, s.key(purpose, email), code, s.ttl).Err(); err != nil {
		return "", err
	}
	return code, nil
}

// Peek checks the supplied code WITHOUT consuming it. Used for a pre-validation
// step (e.g. FE verifying the OTP before showing the new-password form); the
// final consume happens in Verify.
func (s *OTPService) Peek(ctx context.Context, purpose, email, code string) (bool, error) {
	if !s.Enabled() {
		return false, fmt.Errorf("otp backend (redis) not configured")
	}
	stored, err := s.redis.Get(ctx, s.key(purpose, email)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

// Verify checks the supplied code and, on success, deletes it (single use).
func (s *OTPService) Verify(ctx context.Context, purpose, email, code string) (bool, error) {
	if !s.Enabled() {
		return false, fmt.Errorf("otp backend (redis) not configured")
	}
	stored, err := s.redis.Get(ctx, s.key(purpose, email)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if stored != code {
		return false, nil
	}
	_ = s.redis.Del(ctx, s.key(purpose, email)).Err()
	return true, nil
}
