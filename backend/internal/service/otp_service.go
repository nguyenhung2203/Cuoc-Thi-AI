package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type OTPRecord struct {
	Code         string `json:"code"`
	GenerationID string `json:"generation_id"`
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

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
	return fmt.Sprintf("otp:%s:%s", purpose, NormalizeEmail(email))
}

func (s *OTPService) rateKey(purpose, email string) string {
	return fmt.Sprintf("otp:send:rate:%s:%s", purpose, NormalizeEmail(email))
}

func (s *OTPService) AllowSend(ctx context.Context, purpose, email string) (bool, error) {
	if !s.Enabled() {
		return false, fmt.Errorf("otp backend (redis) not configured")
	}
	return s.redis.SetNX(ctx, s.rateKey(purpose, email), "1", time.Minute).Result()
}

// and returns the raw code so the caller can email it.
func (s *OTPService) Generate(ctx context.Context, purpose, email string) (string, string, error) {
	if !s.Enabled() {
		return "", "", fmt.Errorf("otp backend (redis) not configured")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())
	record := OTPRecord{Code: code, GenerationID: uuid.NewString()}
	b, err := json.Marshal(record)
	if err != nil {
		return "", "", err
	}
	if err := s.redis.Set(ctx, s.key(purpose, email), b, s.ttl).Err(); err != nil {
		return "", "", err
	}
	return record.Code, record.GenerationID, nil
}

func (s *OTPService) Current(ctx context.Context, purpose, email string) (code, generationID string, err error) {
	var record OTPRecord
	value, err := s.redis.Get(ctx, s.key(purpose, email)).Result()
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return "", "", err
	}
	return record.Code, record.GenerationID, nil
}

func (s *OTPService) Rollback(ctx context.Context, purpose, email, generationID string) error {
	_, currentGeneration, err := s.Current(ctx, purpose, email)
	if err != nil || currentGeneration != generationID {
		return err
	}
	return s.redis.Del(ctx, s.key(purpose, email), s.rateKey(purpose, email)).Err()
}

// step (e.g. FE verifying the OTP before showing the new-password form); the
// final consume happens in Verify.
func (s *OTPService) Peek(ctx context.Context, purpose, email, code string) (bool, error) {
	if !s.Enabled() {
		return false, fmt.Errorf("otp backend (redis) not configured")
	}
	value, err := s.redis.Get(ctx, s.key(purpose, email)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var record OTPRecord
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return false, err
	}
	return record.Code == code, nil
}

// Verify checks the supplied code and, on success, deletes it (single use).
func (s *OTPService) Verify(ctx context.Context, purpose, email, code string) (bool, error) {
	if !s.Enabled() {
		return false, fmt.Errorf("otp backend (redis) not configured")
	}
	value, err := s.redis.Get(ctx, s.key(purpose, email)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var record OTPRecord
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return false, err
	}
	if record.Code != code {
		return false, nil
	}
	_ = s.redis.Del(ctx, s.key(purpose, email)).Err()
	return true, nil
}
