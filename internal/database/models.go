// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderEMAIL    AuthProvider = "EMAIL"
	AuthProviderGOOGLE   AuthProvider = "GOOGLE"
	AuthProviderFACEBOOK AuthProvider = "FACEBOOK"
	AuthProviderAPPLE    AuthProvider = "APPLE"
)

func (e *AuthProvider) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AuthProvider(s)
	case string:
		*e = AuthProvider(s)
	default:
		return fmt.Errorf("unsupported scan type for AuthProvider: %T", src)
	}
	return nil
}

type NullAuthProvider struct {
	AuthProvider AuthProvider
	Valid        bool // Valid is true if AuthProvider is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAuthProvider) Scan(value interface{}) error {
	if value == nil {
		ns.AuthProvider, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AuthProvider.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAuthProvider) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AuthProvider), nil
}

type ProcessStatus string

const (
	ProcessStatusPENDING    ProcessStatus = "PENDING"
	ProcessStatusPROCESSING ProcessStatus = "PROCESSING"
	ProcessStatusCOMPLETED  ProcessStatus = "COMPLETED"
	ProcessStatusFAILED     ProcessStatus = "FAILED"
)

func (e *ProcessStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProcessStatus(s)
	case string:
		*e = ProcessStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ProcessStatus: %T", src)
	}
	return nil
}

type NullProcessStatus struct {
	ProcessStatus ProcessStatus
	Valid         bool // Valid is true if ProcessStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProcessStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ProcessStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProcessStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProcessStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProcessStatus), nil
}

type Artist struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Name      string
	ImagePath string
	SpotifyID sql.NullString
	ImageUrl  sql.NullString
}

type ArtistAudio struct {
	ID        uuid.UUID
	CreatedAt time.Time
	ArtistID  uuid.UUID
	AudioID   uuid.UUID
}

type Audio struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Title          sql.NullString
	Author         sql.NullString
	DurationMs     sql.NullInt32
	Path           sql.NullString
	SizeBytes      sql.NullInt64
	YoutubeVideoID sql.NullString
	ThumbnailPath  sql.NullString
	SpotifyID      sql.NullString
	ThumbnailUrl   sql.NullString
	LocalID        sql.NullString
}

type AudioLike struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UserID    uuid.UUID
	AudioID   uuid.UUID
}

type Playlist struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	Name              string
	ThumbnailPath     sql.NullString
	SpotifyID         sql.NullString
	ThumbnailUrl      sql.NullString
	AudioImportStatus ProcessStatus
	AudioCount        int32
	TotalAudioCount   int32
}

type PlaylistAudio struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	PlaylistID uuid.UUID
	AudioID    uuid.UUID
}

type User struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Name         sql.NullString
	Email        sql.NullString
	AuthProvider AuthProvider
	PasswordHash sql.NullString
}

type UserAudio struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UserID    uuid.UUID
	AudioID   uuid.UUID
}

type UserPlaylist struct {
	ID                     uuid.UUID
	CreatedAt              time.Time
	UserID                 uuid.UUID
	PlaylistID             uuid.UUID
	IsSpotifySavedPlaylist bool
}

type UserSyncDatum struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	SpotifyLastSyncedAt   sql.NullTime
	UserAudioLastSyncedAt sql.NullTime
}
