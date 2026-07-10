package repository

import (
	"context"
	"time"

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

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newHash, userID)
	return err
}

// Password Reset Token methods
func (r *UserRepository) CreatePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

type PasswordResetTokenRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
}

func (r *UserRepository) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetTokenRow, error) {
	var row PasswordResetTokenRow
	err := r.db.GetContext(ctx, &row, "SELECT id, user_id, token, expires_at, used FROM password_reset_tokens WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenID string) error {
	query := `UPDATE password_reset_tokens SET used = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

// OTP Verifications methods
func (r *UserRepository) CreateOTP(ctx context.Context, email, otp, purpose string, expiresAt time.Time) error {
	query := `
		INSERT INTO otp_verifications (email, otp, purpose, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, email, otp, purpose, expiresAt)
	return err
}

type OTPRow struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	OTP       string    `db:"otp"`
	Purpose   string    `db:"purpose"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
}

func (r *UserRepository) GetOTP(ctx context.Context, email, otp, purpose string) (*OTPRow, error) {
	var row OTPRow
	err := r.db.GetContext(ctx, &row, "SELECT id, email, otp, purpose, expires_at, used FROM otp_verifications WHERE email = $1 AND otp = $2 AND purpose = $3 ORDER BY created_at DESC LIMIT 1", email, otp, purpose)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) MarkOTPUsed(ctx context.Context, id string) error {
	query := `UPDATE otp_verifications SET used = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, userID string) error {
	query := `UPDATE users SET email_verified_at = NOW(), status = 'active' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, email, password_hash, full_name, role, status, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.StructScan(&u); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID, role string) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, role, userID)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, isActive bool) error {
	status := "inactive"
	if isActive {
		status = "active"
	}
	query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, userID)
	return err
}
