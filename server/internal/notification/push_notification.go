package notification

import (
	"context"
	"encoding/json"
	"log/slog"
	"server/internal/config"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
)

type Notification interface {
	// Notify pushes a message to every device the user has subscribed with.
	Notify(ctx context.Context, userId uuid.UUID, title string, body string)
}

type NotificationPayload struct {
	Title string `json:"title"` // must be exported + tagged
	Body  string `json:"body"`
}

type WebPushNotification struct {
	vapidSubject         string
	vapidPublic          string
	vapidPrivate         string
	notificationServices *NotificationService
}

func NewWebPushNotification(cfg config.WebPushNotificationConfig, notificationServices *NotificationService) *WebPushNotification {
	return &WebPushNotification{
		vapidSubject:         cfg.VapidSubject,
		vapidPublic:          cfg.VapidPublic,
		vapidPrivate:         cfg.VapidPrivate,
		notificationServices: notificationServices,
	}
}

func (n *WebPushNotification) Notify(ctx context.Context, userId uuid.UUID, title string, body string) {

	notificationSubscriptions, err := n.notificationServices.GetSubscriptionsOfUser(ctx, userId)
	if err != nil {
		slog.Error("Cannot get subscriptions of user", "userId", userId.String())
		return
	}
	for _, subscription := range notificationSubscriptions {
		payload := NotificationPayload{
			Title: title,
			Body:  body,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			slog.Error("Cannot get subscriptions of user", "userId", userId.String())
			continue
		}
		subs := webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				Auth:   subscription.Auth,
				P256dh: subscription.P256dh,
			},
		}
		resp, err := webpush.SendNotification(payloadBytes, &subs, &webpush.Options{
			Subscriber:      n.vapidSubject, // "mailto:you@..."
			VAPIDPublicKey:  n.vapidPublic,
			VAPIDPrivateKey: n.vapidPrivate,
			TTL:             30,
		})
		if err != nil {
			// On error resp is nil, so guard before touching resp.Body. Close
			// in-loop (not deferred) so bodies don't pile up across iterations.
			slog.Error("failed to send push notification", "error", err.Error(), "endpoint", subscription.Endpoint)
			continue
		}
		resp.Body.Close()
	}
}
