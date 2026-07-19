package notification

import (
	"server/internal/config"
)

type Notification interface {
	Notify(string, string)
}

type NotificationPayload struct {
	Title string `json:"title"` // must be exported + tagged
	Body  string `json:"body"`
}

type WebPushNotification struct {
	vapidSubject string
	vapidPublic  string
	vapidPrivate string
}

func NewWebPushNotification(cfg config.WebPushNotificationConfig) *WebPushNotification {
	return &WebPushNotification{
		vapidSubject: cfg.VapidSubject,
		vapidPublic:  cfg.VapidPublic,
		vapidPrivate: cfg.VapidPrivate,
	}
}

func (n *WebPushNotification) Notify(title string, body string) {
	// payload := NotificationPayload{
	// 	Title: title,
	// 	Body:  body,
	// }
	// payloadBytes, err := json.Marshal(payload)
	// resp, err := webpush.SendNotification(payloadBytes, &sub, &webpush.Options{
	// 	Subscriber:      n.vapidSubject, // "mailto:you@..."
	// 	VAPIDPublicKey:  n.vapidPublic,
	// 	VAPIDPrivateKey: n.vapidPrivate,
	// 	TTL:             30,
	// })
	// defer resp.Body.Close()
}
