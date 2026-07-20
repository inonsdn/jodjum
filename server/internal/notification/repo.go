package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{
		db: db,
	}
}

type NotificationSubscription struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSubscription stores a browser push subscription. The push endpoint is
// the natural unique key: a browser can re-subscribe (e.g. after a key rotation)
// and we just refresh the stored keys instead of creating a duplicate row.
// Requires a UNIQUE constraint on endpoint (see the note when this was added).
func (r *NotificationRepo) CreateSubscription(ctx context.Context, userId uuid.UUID, endpoint string, p256dh string, auth string) (NotificationSubscription, error) {
	const statement = `
		INSERT INTO notification_subscriptions (user_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint) DO UPDATE
		SET user_id = EXCLUDED.user_id, p256dh = EXCLUDED.p256dh, auth = EXCLUDED.auth
		RETURNING id, user_id, endpoint, p256dh, auth, created_at
	`

	var sub NotificationSubscription
	err := r.db.QueryRow(ctx, statement, userId, endpoint, p256dh, auth).Scan(
		&sub.Id,
		&sub.UserId,
		&sub.Endpoint,
		&sub.P256dh,
		&sub.Auth,
		&sub.CreatedAt,
	)
	if err != nil {
		return NotificationSubscription{}, err
	}

	return sub, nil
}

// DeleteSubscriptionByEndpoint removes one of the user's subscriptions. It is
// idempotent: deleting an endpoint that no longer exists is not an error.
func (r *NotificationRepo) DeleteSubscriptionByEndpoint(ctx context.Context, userId uuid.UUID, endpoint string) error {
	const statement = `
		DELETE FROM notification_subscriptions
		WHERE endpoint = $1 AND user_id = $2
	`

	_, err := r.db.Exec(ctx, statement, endpoint, userId)
	return err
}

// GetSubscriptionsOfUser returns every device/browser the user has subscribed
// with, so the reminder loop can push a notification to all of them.
func (r *NotificationRepo) GetSubscriptionsOfUser(ctx context.Context, userId uuid.UUID) ([]NotificationSubscription, error) {
	const statement = `
		SELECT
			id,
			user_id,
			endpoint,
			p256dh,
			auth,
			created_at
		FROM notification_subscriptions
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, statement, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]NotificationSubscription, 0)

	for rows.Next() {
		var sub NotificationSubscription

		err := rows.Scan(
			&sub.Id,
			&sub.UserId,
			&sub.Endpoint,
			&sub.P256dh,
			&sub.Auth,
			&sub.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
