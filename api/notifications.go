package api

import (
	"context"

	"github.com/Uspacy/uspacy-go-sdk/v2/notifications"
)

// GetNotifications returns the notifications.
func (us *Uspacy) GetNotifications(ctx context.Context, opts ...RequestOption) (notifications.Notifications, error) {
	body, err := us.doGet(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), opts...)
	return decodeJSON[notifications.Notifications](body, err)
}

// CreateNotification creates a notification and returns the response with the HTTP status
// code.
func (us *Uspacy) CreateNotification(ctx context.Context, request notifications.CreateNotificationRequest, opts ...RequestOption) (notifications.Notifications, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request, opts...)
	v, err := decodeJSON[notifications.Notifications](body, err)
	return v, code, err
}

// MarkNotificationsAsRead marks notifications as read and returns the HTTP status code.
func (us *Uspacy) MarkNotificationsAsRead(ctx context.Context, request notifications.MarkNotificationsAsReadRequest, opts ...RequestOption) (statusCode int, err error) {
	_, statusCode, err = us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL, notifications.MarkAsReadURL), request, opts...)
	return statusCode, err
}

// DeleteNotification deletes notifications and returns the HTTP status code.
func (us *Uspacy) DeleteNotification(ctx context.Context, request notifications.DeleteNotificationRequest, opts ...RequestOption) (statusCode int, err error) {
	return us.doDelete(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request, opts...)
}
