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

// ListPendingUsers returns all users with status = 'pending'.
func (r *UserRepository) ListPendingUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users WHERE status = 'pending' AND deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ListAllUsers returns all users.
func (r *UserRepository) ListAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateProfile updates editable profile fields for a user.
func (r *UserRepository) UpdateProfile(ctx context.Context, id, fullName, avatarURL string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET full_name = $1, avatar_url = NULLIF($2, ''), updated_at = NOW()
		 WHERE id = $3 AND deleted_at IS NULL`,
		fullName, avatarURL, id)
	return err
}

// UpdateVerificationFile updates the verification_file_id for a user.
func (r *UserRepository) UpdateVerificationFile(ctx context.Context, id, fileID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET verification_file_id = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		fileID, id)
	return err
}

// UpdateStatus updates the status of a user (e.g. pending -> active).
func (r *UserRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET status = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		status, id)
	return err
}

// SoftDelete marks a user account as deleted (data retained, login disabled).
func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET deleted_at = NOW(), status = 'inactive', updated_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`,
		id)
	return err
}

// UpdatePassword sets a new password hash for a user.
func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW()
		 WHERE id = $2 AND deleted_at IS NULL`,
		passwordHash, id)
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

// GetDashboardStats returns high-level metrics for the admin dashboard.
func (r *UserRepository) GetDashboardStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	var totalUsers int
	if err := r.db.GetContext(ctx, &totalUsers, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"); err == nil {
		stats["total_users"] = totalUsers
	}

	var pendingUsers int
	if err := r.db.GetContext(ctx, &pendingUsers, "SELECT COUNT(*) FROM users WHERE status = 'pending' AND deleted_at IS NULL"); err == nil {
		stats["pending_users"] = pendingUsers
	}

	var totalCompanies int
	if err := r.db.GetContext(ctx, &totalCompanies, "SELECT COUNT(*) FROM companies WHERE deleted_at IS NULL"); err == nil {
		stats["total_companies"] = totalCompanies
	}

	var totalInterviews int
	if err := r.db.GetContext(ctx, &totalInterviews, "SELECT COUNT(*) FROM interviews WHERE deleted_at IS NULL"); err == nil {
		stats["total_interviews"] = totalInterviews
	}

	var totalCandidates int
	if err := r.db.GetContext(ctx, &totalCandidates, "SELECT COUNT(*) FROM candidates WHERE deleted_at IS NULL"); err == nil {
		stats["total_candidates"] = totalCandidates
	}

	return stats, nil
}
