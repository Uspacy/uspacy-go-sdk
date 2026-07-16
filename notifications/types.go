package notifications

const (
	VersionURL       = "notifications/v1"
	NotificationsURL = "notifications"
	MarkAsReadURL    = "mark-as-read"
)

type Action string

const (
	ActionCreate      Action = "create"
	ActionUpdate      Action = "update"
	ActionDelete      Action = "delete"
	ActionUpdateStage Action = "updateStage"
)

type NotificationAction string

const (
	NotificationActionSubscribe NotificationAction = "subscribe"
	NotificationActionPublish   NotificationAction = "publish"
)

type Notifications []Notification

type CreateNotificationRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Recipients  []int64 `json:"recipients"`
	Service     string  `json:"service,omitempty"`
	Icon        string  `json:"icon,omitempty"`
	EntityType  string  `json:"entityType,omitempty"`
	EntityID    *int64  `json:"entityId,omitempty"`
}

type ReadNotificationFilter struct {
	CreatedAt int64 `json:"createdAt"`
}

type MarkNotificationsAsReadRequest struct {
	Filter []ReadNotificationFilter `json:"filter"`
}

type DeleteNotificationRequest struct {
	CreatedAt int64 `json:"createdAt"`
}

type Notification struct {
	ID        string           `json:"id"`
	Topic     string           `json:"topic"`
	Type      string           `json:"type"`
	Env       string           `json:"env"`
	Read      *bool            `json:"read,omitempty"`
	CreatedAt int64            `json:"createdAt"`
	Data      NotificationData `json:"data"`
	Metadata  []any            `json:"metadata,omitempty"`
}

type NotificationData struct {
	Entity     any        `json:"entity"`
	RootParent RootParent `json:"root_parent"`
	UserID     int64      `json:"user_id"`
	Service    string     `json:"service"`
	Timestamp  string     `json:"timestamp"`
	Action     Action     `json:"action"`
	Show       bool       `json:"show"`
	OldEntity  any        `json:"old_entity,omitempty"`
}

type RootParent struct {
	Data      RootParentData `json:"data"`
	Type      string         `json:"type"`
	Service   string         `json:"service"`
	TableName string         `json:"table_name"`
}

type RootParentData struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}
