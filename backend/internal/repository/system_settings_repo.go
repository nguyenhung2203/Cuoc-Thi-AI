package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SystemSettingsRepository struct {
	db *sqlx.DB
}

type SystemSettingsMap map[string]interface{}

func NewSystemSettingsRepository(db *sqlx.DB) *SystemSettingsRepository {
	r := &SystemSettingsRepository{db: db}
	if db != nil {
		r.ensureTableAndDefault()
	}
	return r
}

func (r *SystemSettingsRepository) ensureTableAndDefault() {
	query := `
		CREATE TABLE IF NOT EXISTS system_settings (
			config_key VARCHAR(100) PRIMARY KEY,
			config_value JSONB NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		fmt.Printf("ensureTableAndDefault error: %v\n", err)
		return
	}

	defaults := map[string]interface{}{
		"system_name":                    "WeMake AI Recruitment",
		"maintenance_mode":               false,
		"max_upload_size_mb":             20,
		"default_passing_score":          70,
		"default_ai_model":               "gemini-2.5-flash",
		"jwt_token_expiry_hours":         24,
		"admin_2fa_required":             false,
		"notification_ttl_days":          30,
		"notification_max_per_user":      200,
		"enable_email_notifications":     true,
		"notify_on_new_applicant":        true,
		"notify_on_report_ready":         true,
		"notify_on_interview_cancelled":  true,
	}
	for k, v := range defaults {
		b, _ := json.Marshal(v)
		_, _ = r.db.Exec("INSERT INTO system_settings (config_key, config_value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT DO NOTHING", k, b)
	}
}

func (r *SystemSettingsRepository) GetSettingInt(ctx context.Context, key string, defaultVal int) int {
	if r == nil || r.db == nil {
		return defaultVal
	}
	var valBytes []byte
	err := r.db.GetContext(ctx, &valBytes, "SELECT config_value FROM system_settings WHERE config_key = $1", key)
	if err != nil {
		return defaultVal
	}
	var intVal int
	if err := json.Unmarshal(valBytes, &intVal); err != nil {
		var floatVal float64
		if err := json.Unmarshal(valBytes, &floatVal); err == nil {
			return int(floatVal)
		}
		return defaultVal
	}
	return intVal
}

func (r *SystemSettingsRepository) GetSettingBool(ctx context.Context, key string, defaultVal bool) bool {
	if r == nil || r.db == nil {
		return defaultVal
	}
	var valBytes []byte
	err := r.db.GetContext(ctx, &valBytes, "SELECT config_value FROM system_settings WHERE config_key = $1", key)
	if err != nil {
		return defaultVal
	}
	var boolVal bool
	if err := json.Unmarshal(valBytes, &boolVal); err == nil {
		return boolVal
	}
	return defaultVal
}

func (r *SystemSettingsRepository) GetAll(ctx context.Context) (map[string]interface{}, error) {
	type Row struct {
		Key   string `db:"config_key"`
		Value []byte `db:"config_value"`
	}
	var rows []Row
	err := r.db.SelectContext(ctx, &rows, "SELECT config_key, config_value FROM system_settings")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, row := range rows {
		var val interface{}
		if err := json.Unmarshal(row.Value, &val); err == nil {
			result[row.Key] = val
		}
	}
	return result, nil
}

func (r *SystemSettingsRepository) UpdateSettings(ctx context.Context, settings map[string]interface{}) error {
	for k, v := range settings {
		b, err := json.Marshal(v)
		if err != nil {
			continue
		}
		query := `
			INSERT INTO system_settings (config_key, config_value, updated_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (config_key) DO UPDATE
			SET config_value = EXCLUDED.config_value, updated_at = NOW();
		`
		if _, err := r.db.ExecContext(ctx, query, k, b); err != nil {
			return err
		}
	}
	return nil
}
