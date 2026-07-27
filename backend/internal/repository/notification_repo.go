package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type NotificationRepository struct {
	client       *redis.Client
	settingsRepo *SystemSettingsRepository
}

func NewNotificationRepository(client *redis.Client, settingsRepo *SystemSettingsRepository) *NotificationRepository {
	return &NotificationRepository{client: client, settingsRepo: settingsRepo}
}

func (r *NotificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	if r.settingsRepo != nil {
		switch notif.Type {
		case "new_applicant":
			if !r.settingsRepo.GetSettingBool(ctx, "notify_on_new_applicant", true) {
				return nil
			}
		case "report_ready", "report_failed":
			if !r.settingsRepo.GetSettingBool(ctx, "notify_on_report_ready", true) {
				return nil
			}
		case "interview_cancelled":
			if !r.settingsRepo.GetSettingBool(ctx, "notify_on_interview_cancelled", true) {
				return nil
			}
		}
	}

	id, err := r.client.Incr(ctx, "notif:next_id").Result()
	if err != nil {
		return err
	}

	notif.ID = uint64(id)
	if notif.CreatedAt.IsZero() {
		notif.CreatedAt = time.Now()
	}

	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	maxPerUser := 200
	ttlDays := 30
	if r.settingsRepo != nil {
		maxPerUser = r.settingsRepo.GetSettingInt(ctx, "notification_max_per_user", 200)
		if maxPerUser <= 0 {
			maxPerUser = 200
		}
		ttlDays = r.settingsRepo.GetSettingInt(ctx, "notification_ttl_days", 30)
		if ttlDays <= 0 {
			ttlDays = 30
		}
	}

	userKey := fmt.Sprintf("notif:user:%s:list", notif.UserID)
	pipe := r.client.Pipeline()
	pipe.LPush(ctx, userKey, data)
	pipe.LTrim(ctx, userKey, 0, int64(maxPerUser-1))
	pipe.Expire(ctx, userKey, time.Duration(ttlDays)*24*time.Hour)
	pipe.Set(ctx, fmt.Sprintf("notif:id_to_user:%d", notif.ID), notif.UserID, time.Duration(ttlDays)*24*time.Hour)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	_ = r.client.Publish(ctx, fmt.Sprintf("notif:pubsub:user:%s", notif.UserID), data).Err()
	return nil
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error) {
	if r.client == nil {
		return []models.Notification{}, nil
	}

	userKey := fmt.Sprintf("notif:user:%s:list", userID)
	start := int64(offset)
	stop := int64(offset + limit - 1)
	if limit <= 0 {
		stop = start + 19
	}

	rawItems, err := r.client.LRange(ctx, userKey, start, stop).Result()
	if err != nil {
		return nil, err
	}

	items := make([]models.Notification, 0, len(rawItems))
	for _, raw := range rawItems {
		var notif models.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err == nil {
			items = append(items, notif)
		}
	}
	return items, nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id uint64, userID string) error {
	if r.client == nil {
		return sql.ErrNoRows
	}

	userKey := fmt.Sprintf("notif:user:%s:list", userID)
	rawItems, err := r.client.LRange(ctx, userKey, 0, -1).Result()
	if err != nil {
		return err
	}

	found := false
	for i, raw := range rawItems {
		var notif models.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err == nil {
			if notif.ID == id {
				found = true
				if !notif.IsRead {
					notif.IsRead = true
					newData, _ := json.Marshal(notif)
					if err := r.client.LSet(ctx, userKey, int64(i), newData).Err(); err != nil {
						return err
					}
				}
				break
			}
		}
	}

	if !found {
		return sql.ErrNoRows
	}
	return nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	if r.client == nil {
		return sql.ErrNoRows
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	userID, err := r.client.Get(ctx, fmt.Sprintf("notif:id_to_user:%s", id)).Result()
	if err != nil || userID == "" {
		return sql.ErrNoRows
	}
	return r.MarkRead(ctx, idUint, userID)
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	if r.client == nil {
		return nil
	}

	userKey := fmt.Sprintf("notif:user:%s:list", userID)
	rawItems, err := r.client.LRange(ctx, userKey, 0, -1).Result()
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()
	modified := false
	for i, raw := range rawItems {
		var notif models.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err == nil {
			if !notif.IsRead {
				notif.IsRead = true
				newData, _ := json.Marshal(notif)
				pipe.LSet(ctx, userKey, int64(i), newData)
				modified = true
			}
		}
	}

	if modified {
		_, err = pipe.Exec(ctx)
		return err
	}
	return nil
}

func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID string) (int, error) {
	if r.client == nil {
		return 0, nil
	}

	userKey := fmt.Sprintf("notif:user:%s:list", userID)
	rawItems, err := r.client.LRange(ctx, userKey, 0, -1).Result()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, raw := range rawItems {
		var notif models.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err == nil {
			if !notif.IsRead {
				count++
			}
		}
	}
	return count, nil
}

