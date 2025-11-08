package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	PublicVisibility   = "public"
	PrivateVisibility  = "private"
	unlistedVisibility = "unlisted"
)

type Share struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	FromUserID  uuid.UUID `gorm:"type:uuid;index" json:"from_user_id"`
	ToUserID    uuid.UUID `gorm:"type:uuid;index" json:"to_user_id"`
	WatchlistID uuid.UUID `gorm:"type:uuid;index" json:"watchlist_id"`
	Message     string    `json:"message"`
}

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"` // recipient
	Type      string    `gorm:"type:text;not null;check:type IN ('like','follow')"`
	ActorID   uuid.UUID `gorm:"type:uuid;not null"`
	EntityID  uuid.UUID `gorm:"type:uuid;not null"`
	IsRead    bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

type Activity struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Type      string    `gorm:"type:text;not null;check:type IN ('like','follow','create_list','add_item')"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (Activity) TableName() string { return "activities" }
func EncodeStringSlice(ss []string) datatypes.JSON {
	if ss == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(ss)
	return datatypes.JSON(b)
}

func DecodeStringSlice(j datatypes.JSON) []string {
	var out []string
	_ = json.Unmarshal(j, &out)
	return out
}
