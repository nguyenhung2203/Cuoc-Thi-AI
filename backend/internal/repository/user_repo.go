package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"backend/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, status)
		VALUES (:id, :email, :password_hash, :full_name, :role, :status)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates editable profile fields for a user.
func (r *UserRepository) UpdateProfile(ctx context.Context, id, fullName, avatarURL string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET full_name = $1, avatar_url = NULLIF($2, ''), updated_at = NOW()
		 WHERE id = $3 AND deleted_at IS NULL`,
		fullName, avatarURL, id)
	return err
}

// UpdateSettings replaces the user's settings JSON blob.
func (r *UserRepository) UpdateSettings(ctx context.Context, id string, settings models.JSONB) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET settings = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		settings, id)
	return err
}

type UserCompanyRow struct {
	CompanyID   string `db:"company_id"`
	CompanyName string `db:"company_name"`
	Role        string `db:"role"`
}

func (r *UserRepository) FindUserCompanies(ctx context.Context, userID string) ([]UserCompanyRow, error) {
	query := `
		SELECT c.id AS company_id, c.name AS company_name, cm.role
		FROM company_members cm
		JOIN companies c ON c.id = cm.company_id
		WHERE cm.user_id = $1 AND cm.status = 'active' AND c.deleted_at IS NULL
	`
	var rows []UserCompanyRow
	err := r.db.SelectContext(ctx, &rows, query, userID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
