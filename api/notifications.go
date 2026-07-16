package api

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/notifications"
)

func (us *Uspacy) GetNotifications() (entities notifications.Notifications, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(notifications.VersionURL, notifications.NotificationsURL))
	if err != nil {
		return entities, err
	}

	err = json.Unmarshal(body, &entities)
	return entities, err
}

func (us *Uspacy) CreateNotification(request notifications.CreateNotificationRequest) (entities notifications.Notifications, statusCode int, err error) {
	body, statusCode, err := us.doPost(us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
	if err != nil {
		return entities, statusCode, err
	}

	err = json.Unmarshal(body, &entities)
	return entities, statusCode, err
}

func (us *Uspacy) MarkNotificationsAsRead(request notifications.MarkNotificationsAsReadRequest) (statusCode int, err error) {
	_, statusCode, err = us.doPost(us.buildURL(notifications.VersionURL, notifications.NotificationsURL, notifications.MarkAsReadURL), request)
	return statusCode, err
}

func (us *Uspacy) DeleteNotification(request notifications.DeleteNotificationRequest) (statusCode int, err error) {
	return us.doDeleteEmptyHeaders(us.buildURL(notifications.VersionURL, notifications.NotificationsURL), request)
}
