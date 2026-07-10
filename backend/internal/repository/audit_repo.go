package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type AuditRepository struct {
	db *sqlx.DB
}

func NewAuditRepository(db *sqlx.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Insert(ctx context.Context, al *models.AuditLog) error {
	q := `
		INSERT INTO audit_logs (
			id, company_id, actor_user_id, actor_role,
			action, resource_type, resource_id,
			before_json, after_json, ip_address, user_agent,
			created_at
		) VALUES (
			:id, :company_id, :actor_user_id, :actor_role,
			:action, :resource_type, :resource_id,
			:before_json, :after_json, :ip_address, :user_agent,
			NOW()
		) RETURNING id, created_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowContext(ctx, al).Scan(&al.ID, &al.CreatedAt)
}

func (r *AuditRepository) ListByCompany(ctx context.Context, companyID string, resourceType, actorUserID string, limit, offset int) ([]models.AuditLog, error) {
	args := []interface{}{companyID}
	query := `SELECT * FROM audit_logs WHERE company_id = $1`
	paramIdx := 2

	if resourceType != "" {
		query += ` AND resource_type = $` + string(rune('0'+paramIdx))
		args = append(args, resourceType)
		paramIdx++
	}
	if actorUserID != "" {
		query += ` AND actor_user_id = $` + string(rune('0'+paramIdx))
		args = append(args, actorUserID)
		paramIdx++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+paramIdx)) + ` OFFSET $` + string(rune('0'+paramIdx+1))
	args = append(args, limit, offset)

	var items []models.AuditLog
	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuditRepository) ListSystem(ctx context.Context, limit, offset int) ([]models.AuditLog, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE action NOT IN ('login', 'logout')`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT * FROM audit_logs WHERE action NOT IN ('login', 'logout') ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	var items []models.AuditLog
	if err := r.db.SelectContext(ctx, &items, query, limit, offset); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
