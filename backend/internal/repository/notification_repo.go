package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"backend/internal/models"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	q := `
		INSERT INTO notifications (user_id, title, message, type, is_read, link, created_at)
		VALUES (:user_id, :title, :message, :type, :is_read, :link, NOW())
		RETURNING id, created_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, notif).Scan(&notif.ID, &notif.CreatedAt)
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error) {
	q := `SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var items []models.Notification
	err := r.db.SelectContext(ctx, &items, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id uint64, userID string) error {
	q := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	q := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	q := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID string) (int, error) {
	q := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.GetContext(ctx, &count, q, userID)
	return count, err
}
