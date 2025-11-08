package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourname/moodle/internal/models"
)

// Share represents a watchlist share in the domain layer
type Share struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	FromUserID  uuid.UUID
	ToUserID    uuid.UUID
	WatchlistID uuid.UUID
	Message     string
}

// FromModel converts models.Share to domain.Share
func (s *Share) FromModel(model *models.Share) *Share {
	if model == nil {
		return nil
	}
	return &Share{
		ID:          model.ID,
		CreatedAt:   model.CreatedAt,
		FromUserID:  model.FromUserID,
		ToUserID:    model.ToUserID,
		WatchlistID: model.WatchlistID,
		Message:     model.Message,
	}
}

// ToModel converts domain.Share to models.Share
func (s *Share) ToModel() *models.Share {
	if s == nil {
		return nil
	}
	return &models.Share{
		ID:          s.ID,
		CreatedAt:   s.CreatedAt,
		FromUserID:  s.FromUserID,
		ToUserID:    s.ToUserID,
		WatchlistID: s.WatchlistID,
		Message:     s.Message,
	}
}

// ShareFromModel is a helper function to convert models.Share to domain.Share
func ShareFromModel(model *models.Share) *Share {
	if model == nil {
		return nil
	}
	var s Share
	return s.FromModel(model)
}

// SharesFromModel converts a slice of models.Share to domain.Share
func SharesFromModel(modelShares []models.Share) []Share {
	if modelShares == nil {
		return nil
	}
	shares := make([]Share, 0, len(modelShares))
	for _, m := range modelShares {
		shares = append(shares, *ShareFromModel(&m))
	}
	return shares
}

// Notification represents a notification in the domain layer
type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      string
	ActorID   uuid.UUID
	EntityID  uuid.UUID
	IsRead    bool
	CreatedAt time.Time
}

// FromModel converts models.Notification to domain.Notification
func (n *Notification) FromModel(model *models.Notification) *Notification {
	if model == nil {
		return nil
	}
	return &Notification{
		ID:        model.ID,
		UserID:    model.UserID,
		Type:      model.Type,
		ActorID:   model.ActorID,
		EntityID:  model.EntityID,
		IsRead:    model.IsRead,
		CreatedAt: model.CreatedAt,
	}
}

// ToModel converts domain.Notification to models.Notification
func (n *Notification) ToModel() *models.Notification {
	if n == nil {
		return nil
	}
	return &models.Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		ActorID:   n.ActorID,
		EntityID:  n.EntityID,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

// NotificationFromModel is a helper function to convert models.Notification to domain.Notification
func NotificationFromModel(model *models.Notification) *Notification {
	if model == nil {
		return nil
	}
	var n Notification
	return n.FromModel(model)
}

// NotificationsFromModel converts a slice of models.Notification to domain.Notification
func NotificationsFromModel(modelNotifications []models.Notification) []Notification {
	if modelNotifications == nil {
		return nil
	}
	notifications := make([]Notification, 0, len(modelNotifications))
	for _, m := range modelNotifications {
		notifications = append(notifications, *NotificationFromModel(&m))
	}
	return notifications
}

// Activity represents an activity in the domain layer
type Activity struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      string
	SubjectID uuid.UUID
	CreatedAt time.Time
}

// FromModel converts models.Activity to domain.Activity
func (a *Activity) FromModel(model *models.Activity) *Activity {
	if model == nil {
		return nil
	}
	return &Activity{
		ID:        model.ID,
		UserID:    model.UserID,
		Type:      model.Type,
		SubjectID: model.SubjectID,
		CreatedAt: model.CreatedAt,
	}
}

// ToModel converts domain.Activity to models.Activity
func (a *Activity) ToModel() *models.Activity {
	if a == nil {
		return nil
	}
	return &models.Activity{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      a.Type,
		SubjectID: a.SubjectID,
		CreatedAt: a.CreatedAt,
	}
}

// ActivityFromModel is a helper function to convert models.Activity to domain.Activity
func ActivityFromModel(model *models.Activity) *Activity {
	if model == nil {
		return nil
	}
	var a Activity
	return a.FromModel(model)
}

// ActivitiesFromModel converts a slice of models.Activity to domain.Activity
func ActivitiesFromModel(modelActivities []models.Activity) []Activity {
	if modelActivities == nil {
		return nil
	}
	activities := make([]Activity, 0, len(modelActivities))
	for _, m := range modelActivities {
		activities = append(activities, *ActivityFromModel(&m))
	}
	return activities
}
