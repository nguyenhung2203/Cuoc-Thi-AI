package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// matchCacheTTL bounds staleness of a cached CV↔job fit score. A score is a
// function of (parsed CV, job requirements); both change rarely, so a day is a
// safe default. CV re-upload invalidates explicitly (see InvalidateUser).
const matchCacheTTL = 24 * time.Hour

// MatchCacheService caches AI-computed CV↔job fit results in Redis so the
// browse/preview path doesn't re-run the model on every view. It degrades
// gracefully: with no Redis wired, every method is a no-op miss and callers
// fall back to computing inline.
type MatchCacheService struct {
	redis *redis.Client
}

func NewMatchCacheService(rc *redis.Client) *MatchCacheService {
	return &MatchCacheService{redis: rc}
}

func (s *MatchCacheService) Enabled() bool {
	return s != nil && s.redis != nil
}

func (s *MatchCacheService) key(userID, jobID string) string {
	return fmt.Sprintf("fitmatch:%s:%s", userID, jobID)
}

// Get returns the cached match for (user, job). The second value is false on a
// cache miss (including when Redis is unavailable).
func (s *MatchCacheService) Get(ctx context.Context, userID, jobID string) (*MatchResult, bool) {
	if !s.Enabled() {
		return nil, false
	}
	raw, err := s.redis.Get(ctx, s.key(userID, jobID)).Result()
	if err != nil {
		return nil, false
	}
	var result MatchResult
	if json.Unmarshal([]byte(raw), &result) != nil {
		return nil, false
	}
	return &result, true
}

// Set stores the match for (user, job) with the default TTL. Best-effort:
// failures are swallowed so caching never breaks the request path.
func (s *MatchCacheService) Set(ctx context.Context, userID, jobID string, result *MatchResult) {
	if !s.Enabled() || result == nil {
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		return
	}
	_ = s.redis.Set(ctx, s.key(userID, jobID), b, matchCacheTTL).Err()
}

// InvalidateUser drops all cached matches for a user. Called when the user's CV
// changes, since every score for that user is now stale.
func (s *MatchCacheService) InvalidateUser(ctx context.Context, userID string) {
	if !s.Enabled() {
		return
	}
	pattern := fmt.Sprintf("fitmatch:%s:*", userID)
	iter := s.redis.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		_ = s.redis.Del(ctx, keys...).Err()
	}
}
