package shared

import (
	"log"
	"os/exec"
)

const (
	ErrInternal                      = "INTERNAL"
	ErrNotFound                      = "NOT_FOUND"
	ErrUnauthorized                  = "UNAUTHORIZED"
	ErrInvalidJSON                   = "INVALID_JSON"
	ErrInvalidGoogleToken            = "INVALID_GOOGLE_TOKEN"
	ErrUserNotFound                  = "USER_NOT_FOUND"
	ErrInvalidQueryParams            = "INVALID_QUERY_PARAMS"
	ErrMissingToken                  = "MISSING_TOKEN"
	ErrInvalidToken                  = "INVALID_TOKEN"
	ErrAudioAlreadyExists            = "AUDIO_ALREADY_EXISTS"
	ErrUserAudioAlreadyExists        = "USER_AUDIO_ALREADY_EXISTS"
	ErrPlaylistAudioAlreadyExists    = "PLAYLIST_AUDIO_ALREADY_EXISTS"
	ErrInvalidAuthMethod             = "INVALID_AUTH_METHOD"
	ErrInvalidEmailOrPassword        = "INVALID_EMAIL_OR_PASSWORD"
	ErrMethodNotAllowed              = "METHOD_NOT_ALLOWED"
	ErrExceededMaxUploadSize         = "EXCEEDED_MAX_UPLOAD_SIZE"
	ErrInvalidMimeType               = "INVALID_MIME_TYPE"
	ErrFailedToGetSpotifyPlaylist    = "FAILED_TO_GET_SPOTIFY_PLAYLIST"
	ErrFailedToGetSpotifyAccessToken = "FAILED_TO_GET_SPOTIFY_ACCESS_TOKEN"
	ErrFailedToGetSpotifyPlaylists   = "FAILED_TO_GET_SPOTIFY_PLAYLISTS"
	ErrUserSyncDatumNotFound         = "USER_SYNC_DATUM_NOT_FOUND"
	ErrPlaylistNotFound              = "PLAYLIST_NOT_FOUND"
	ErrAudioLikeNotFound             = "AUDIO_LIKE_NOT_FOUND"
	ErrAudioNotFound                 = "AUDIO_NOT_FOUND"
	ErrNotImplemented                = "NOT_IMPLEMENTED"
	ErrSocketNotFound                = "SOCKET_NOT_FOUND"
	ErrUserAudioNotFound             = "USER_AUDIO_NOT_FOUND"
	ErrInvalidUUID                   = "INVALID_UUID"
	ErrUserPlaylistNotFound          = "USER_PLAYLIST_NOT_FOUND"
	ErrUserPlaylistNotYours          = "USER_PLAYLIST_NOT_YOURS"
)

func LogCommandError(err error, label string) {
	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf("[%s] Command stderr: %s", label, exitError.Stderr)
	} else {
		log.Printf("[%s] Error executing command: %v", label, err)
	}
}
