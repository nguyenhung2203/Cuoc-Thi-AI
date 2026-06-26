package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"backend/internal/models"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	RevokeFamily(ctx context.Context, familyID string) error
	RevokeAllByUserID(ctx context.Context, userID string) (int, error)
	MarkAsRevoked(ctx context.Context, id string) error
}

type refreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, family_id, token_hash, is_revoked, ip_address, user_agent, expires_at)
		VALUES (:user_id, :family_id, :token_hash, :is_revoked, :ip_address, :user_agent, :expires_at)
		RETURNING id, created_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, token)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(token)
		return err
	}
	return nil
}

func (r *refreshTokenRepo) FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `SELECT * FROM refresh_tokens WHERE token_hash = $1`
	err := r.db.GetContext(ctx, &token, query, hash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepo) RevokeFamily(ctx context.Context, familyID string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE family_id = $1`
	_, err := r.db.ExecContext(ctx, query, familyID)
	return err
}

func (r *refreshTokenRepo) RevokeAllByUserID(ctx context.Context, userID string) (int, error) {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE user_id = $1 AND is_revoked = false`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return int(rows), nil
}

func (r *refreshTokenRepo) MarkAsRevoked(ctx context.Context, id string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
