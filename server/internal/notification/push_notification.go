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
	// Confirm VAPID config is loaded without leaking the private key.
	slog.Info("web push configured",
		"subject", cfg.VapidSubject,
		"publicKeySet", cfg.VapidPublic != "",
		"privateKeySet", cfg.VapidPrivate != "")

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
	slog.Info("push: sending to subscriptions", "userId", userId.String(), "count", len(notificationSubscriptions), "title", title)
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
		// A 2xx (usually 201) means the push service accepted it. 4xx/410 means
		// the subscription is invalid/expired and should be pruned.
		if resp.StatusCode >= 400 {
			slog.Warn("push service rejected notification", "status", resp.StatusCode, "endpoint", subscription.Endpoint)
		} else {
			slog.Info("push sent", "status", resp.StatusCode, "endpoint", subscription.Endpoint)
		}
		resp.Body.Close()
	}
}
