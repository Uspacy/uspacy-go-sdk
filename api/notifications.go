package api

import (
	"context"

	"github.com/Uspacy/uspacy-go-sdk/v2/notifications"
)

// GetNotifications returns notifications list
func (us *Uspacy) GetNotifications(ctx context.Context) (notifications.Notifications, error) {
	body, err := us.doGet(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL))
	return decodeJSON[notifications.Notifications](body, err)
}

// CreateNotification creates a notification
func (us *Uspacy) CreateNotification(ctx context.Context, request notifications.CreateNotificationRequest) (notifications.Notifications, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
	v, err := decodeJSON[notifications.Notifications](body, err)
	return v, code, err
}

// MarkNotificationsAsRead marks notifications as read
func (us *Uspacy) MarkNotificationsAsRead(ctx context.Context, request notifications.MarkNotificationsAsReadRequest) (statusCode int, err error) {
	_, statusCode, err = us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL, notifications.MarkAsReadURL), request)
	return statusCode, err
}

// DeleteNotification deletes a notification
func (us *Uspacy) DeleteNotification(ctx context.Context, request notifications.DeleteNotificationRequest) (statusCode int, err error) {
	return us.doDelete(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
}
