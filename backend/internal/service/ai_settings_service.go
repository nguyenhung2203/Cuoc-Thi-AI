package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

const (
	DefaultTextModel  = "gemini-3.1-flash-lite"
	DefaultVoiceModel = "gemini-2.5-flash-native-audio-latest"
	apiKeySetting     = "gemini_api_key_encrypted"
)

var textModels = map[string]bool{
	"gemini-3.1-flash-lite": true,
	"gemini-2.5-flash-lite": true,
	"gemini-2.5-flash":      true,
}

var voiceModels = map[string]bool{
	DefaultVoiceModel: true,
}

type AIRuntimeConfig struct {
	APIKeys   []string `json:"api_keys"`
	TextModel string   `json:"text_model"`
	LiveModel string   `json:"voice_model"`
}

type AIRuntimeConfigProvider interface {
	GetRuntimeConfig(context.Context) (AIRuntimeConfig, error)
}

type AISettingsView struct {
	TextModel        string `json:"text_model"`
	VoiceModel       string `json:"voice_model"`
	APIKeyConfigured bool   `json:"api_key_configured"`
	APIKeyMasked     string `json:"api_key_masked,omitempty"`
}

type UpdateAISettingsInput struct {
	TextModel   string `json:"text_model"`
	VoiceModel  string `json:"voice_model"`
	APIKey      string `json:"api_key"`
	ClearAPIKey bool   `json:"clear_api_key"`
}

type AISettingsService struct {
	repo          *repository.SystemSettingsRepository
	encryptionKey []byte
	envKeys       []string
	envTextModel  string
	envLiveModel  string
}

func NewAISettingsService(repo *repository.SystemSettingsRepository, encryptionSecret, envKeys, envTextModel, envLiveModel string) *AISettingsService {
	var key []byte
	if strings.TrimSpace(encryptionSecret) != "" {
		sum := sha256.Sum256([]byte(encryptionSecret))
		key = sum[:]
	}
	return &AISettingsService{
		repo: repo, encryptionKey: key, envKeys: splitAPIKeys(envKeys),
		envTextModel: envTextModel, envLiveModel: envLiveModel,
	}
}

func NewAISettingsServiceFromEnv(repo *repository.SystemSettingsRepository) *AISettingsService {
	keys := os.Getenv("GEMINI_API_KEYS")
	if keys == "" {
		keys = os.Getenv("GEMINI_API_KEY")
	}
	return NewAISettingsService(repo, os.Getenv("AI_SETTINGS_ENCRYPTION_KEY"), keys, os.Getenv("AI_DEFAULT_MODEL"), os.Getenv("GEMINI_LIVE_MODEL"))
}

func (s *AISettingsService) GetSettings(ctx context.Context) (AISettingsView, error) {
	cfg, err := s.GetRuntimeConfig(ctx)
	if err != nil {
		return AISettingsView{}, errors.NewInternal("failed to retrieve AI settings")
	}
	view := AISettingsView{TextModel: cfg.TextModel, VoiceModel: cfg.LiveModel}
	if len(cfg.APIKeys) > 0 && !isPlaceholderKey(cfg.APIKeys[0]) {
		view.APIKeyConfigured = true
		view.APIKeyMasked = maskAPIKey(cfg.APIKeys[0])
	}
	return view, nil
}

func (s *AISettingsService) UpdateSettings(ctx context.Context, in UpdateAISettingsInput) (AISettingsView, error) {
	in.TextModel = strings.TrimSpace(in.TextModel)
	in.VoiceModel = strings.TrimSpace(in.VoiceModel)
	in.APIKey = strings.TrimSpace(in.APIKey)
	if !textModels[in.TextModel] {
		return AISettingsView{}, errors.NewBadRequest("unsupported Gemini text model")
	}
	if !voiceModels[in.VoiceModel] {
		return AISettingsView{}, errors.NewBadRequest("unsupported Gemini Live model")
	}
	if in.ClearAPIKey && in.APIKey != "" {
		return AISettingsView{}, errors.NewBadRequest("api_key and clear_api_key cannot be used together")
	}
	updates := map[string]interface{}{"default_ai_model": in.TextModel, "default_ai_voice_model": in.VoiceModel}
	if in.APIKey != "" {
		if len(s.encryptionKey) == 0 {
			return AISettingsView{}, errors.NewInternal("AI settings encryption is not configured")
		}
		ciphertext, err := encryptSecret(s.encryptionKey, in.APIKey)
		if err != nil {
			return AISettingsView{}, errors.NewInternal("failed to protect API key")
		}
		updates[apiKeySetting] = ciphertext
	}
	if err := s.repo.UpdateSettings(ctx, updates); err != nil {
		return AISettingsView{}, errors.NewInternal("failed to update AI settings")
	}
	if in.ClearAPIKey {
		if err := s.repo.DeleteSetting(ctx, apiKeySetting); err != nil {
			return AISettingsView{}, errors.NewInternal("failed to clear API key")
		}
		s.envKeys = nil
	}
	return s.GetSettings(ctx)
}

func (s *AISettingsService) GetRuntimeConfig(ctx context.Context) (AIRuntimeConfig, error) {
	textModel := s.repo.GetSettingString(ctx, "default_ai_model", firstNonEmpty(s.envTextModel, DefaultTextModel))
	liveModel := s.repo.GetSettingString(ctx, "default_ai_voice_model", firstNonEmpty(s.envLiveModel, DefaultVoiceModel))
	keys := append([]string(nil), s.envKeys...)
	ciphertext := s.repo.GetSettingString(ctx, apiKeySetting, "")
	if ciphertext != "" {
		if len(s.encryptionKey) == 0 {
			return AIRuntimeConfig{}, fmt.Errorf("AI settings encryption is not configured")
		}
		plain, err := decryptSecret(s.encryptionKey, ciphertext)
		if err != nil {
			return AIRuntimeConfig{}, fmt.Errorf("decrypt API key: %w", err)
		}
		keys = splitAPIKeys(plain)
	}
	return AIRuntimeConfig{APIKeys: keys, TextModel: textModel, LiveModel: liveModel}, nil
}

func splitAPIKeys(value string) []string {
	var result []string
	for _, key := range strings.Split(value, ",") {
		key = strings.TrimSpace(key)
		if key != "" && !isPlaceholderKey(key) {
			result = append(result, key)
		}
	}
	return result
}

func isPlaceholderKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	return key == "" || key == "your-gemini-api-key" || key == "changeme"
}

func maskAPIKey(key string) string {
	if len(key) <= 4 {
		return "••••"
	}
	return "••••••••" + key[len(key)-4:]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func encryptSecret(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func decryptSecret(key []byte, ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("invalid ciphertext")
	}
	plain, err := gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
	return string(plain), err
}
