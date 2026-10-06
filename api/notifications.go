package api

import (
	"context"
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/v2/notifications"
)

// GetNotifications returns notifications list
func (us *Uspacy) GetNotifications(ctx context.Context) (entities notifications.Notifications, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL))
	if err != nil {
		return entities, err
	}

	err = json.Unmarshal(body, &entities)
	return entities, err
}

// CreateNotification creates a notification
func (us *Uspacy) CreateNotification(ctx context.Context, request notifications.CreateNotificationRequest) (entities notifications.Notifications, statusCode int, err error) {
	body, statusCode, err := us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
	if err != nil {
		return entities, statusCode, err
	}

	err = json.Unmarshal(body, &entities)
	return entities, statusCode, err
}

// MarkNotificationsAsRead marks notifications as read
func (us *Uspacy) MarkNotificationsAsRead(ctx context.Context, request notifications.MarkNotificationsAsReadRequest) (statusCode int, err error) {
	_, statusCode, err = us.doPost(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL, notifications.MarkAsReadURL), request)
	return statusCode, err
}

// DeleteNotification deletes a notification
func (us *Uspacy) DeleteNotification(ctx context.Context, request notifications.DeleteNotificationRequest) (statusCode int, err error) {
	return us.doDeleteEmptyHeaders(ctx, us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
}
