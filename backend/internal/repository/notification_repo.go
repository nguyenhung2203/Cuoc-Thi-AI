package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	q := `
		INSERT INTO notifications (
			user_id, type, title, content, data_json, read_at, created_at
		) VALUES (
			:user_id, :type, :title, :content, :data_json, :read_at, :created_at
		) RETURNING id
	`
	rows, err := r.db.NamedQueryContext(ctx, q, n)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&n.ID)
	}
	return nil
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error) {
	q := `SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var items []models.Notification
	err := r.db.SelectContext(ctx, &items, q, userID, limit, offset)
	return items, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	q := `UPDATE notifications SET read_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	q := `UPDATE notifications SET read_at = NOW() WHERE user_id = $1 AND read_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}
