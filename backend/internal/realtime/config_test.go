package realtime

import (
	"testing"

	"backend/internal/livekit"
)

func TestAllowDevBroadcast(t *testing.T) {
	cases := []struct {
		appEnv, flag string
		want         bool
	}{
		// Production is hard-off regardless of the flag.
		{"production", "true", false},
		{"prod", "true", false},
		{"PRODUCTION", "true", false},
		// Dev requires the explicit opt-in.
		{"development", "true", true},
		{"development", "TRUE", true},
		{"development", "", false},
		{"development", "false", false},
		{"", "true", true}, // unset env counts as non-production
		{"", "", false},
	}
	for _, c := range cases {
		if got := allowDevBroadcast(c.appEnv, c.flag); got != c.want {
			t.Errorf("allowDevBroadcast(%q, %q) = %v, want %v", c.appEnv, c.flag, got, c.want)
		}
	}
}

func TestLoadLiveKitConfig_ProdRejectsDevSecret(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("LIVEKIT_API_KEY", "real-key")
	t.Setenv("LIVEKIT_API_SECRET", "devsecret")
	t.Setenv("LIVEKIT_URL", "wss://livekit.real.example")

	if _, err := livekit.Load(); err == nil {
		t.Fatal("production with devsecret must fail fast")
	}

	t.Setenv("LIVEKIT_API_SECRET", "a-real-secret")
	if _, err := livekit.Load(); err != nil {
		t.Fatalf("production with real credentials should load: %v", err)
	}
}

func TestLoadLiveKitConfig_DevFillsFallbacks(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("LIVEKIT_API_KEY", "")
	t.Setenv("LIVEKIT_API_SECRET", "")
	t.Setenv("LIVEKIT_URL", "")

	cfg, err := livekit.Load()
	if err != nil {
		t.Fatalf("dev load: %v", err)
	}
	if cfg.APISecret != livekit.DevSecret || cfg.APIKey != livekit.DevKey {
		t.Fatalf("dev fallbacks not applied: %+v", cfg)
	}
}
