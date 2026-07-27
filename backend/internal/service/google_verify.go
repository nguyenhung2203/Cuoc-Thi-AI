package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GoogleTokenInfo is the subset of Google's tokeninfo response we rely on.
type GoogleTokenInfo struct {
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Exp           string `json:"exp"`
}

// verifyGoogleIDToken validates a Google-issued ID token against Google's
// tokeninfo endpoint and checks the audience matches our configured client ID.
// Returns the verified claims so the caller never trusts client-supplied email.
func verifyGoogleIDToken(ctx context.Context, idToken, clientID string) (*GoogleTokenInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("missing id_token")
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to reach Google token verifier: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid google id_token")
	}

	var info GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to parse google token info: %w", err)
	}

	if clientID != "" && info.Aud != clientID {
		return nil, fmt.Errorf("google id_token audience mismatch")
	}
	if info.Email == "" || info.EmailVerified != "true" {
		return nil, fmt.Errorf("google account email not verified")
	}

	return &info, nil
}
