package audio

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadYoutubeAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[downloadYoutubeAudioDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		// check if the user already has the audio
		if _, err := GetUserAudioByYoutubeVideoId(r.Context(), apiCfg.DB, authPayload.UserId, body.VideoId); err == nil {
			shared.ResConflict(w, shared.ErrAudioAlreadyExists)
			return
		}

		videoInfo, err := youtube.GetYoutubeVideoInfo(body.VideoId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		filePath, thumbnailPath, err := youtube.DownloadYoutubeAudio(body.VideoId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		fileSize, err := shared.GetFileSize(filePath)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		newAudio, err := CreateAudio(
			r.Context(),
			apiCfg.DB,
			database.CreateAudioParams{
				Title:          sql.NullString{String: strings.TrimSpace(videoInfo.Title), Valid: true},
				Author:         sql.NullString{String: strings.TrimSpace(videoInfo.Uploader), Valid: true},
				DurationMs:     sql.NullInt32{Int32: int32(videoInfo.DurationSeconds * 1000), Valid: true},
				Path:           sql.NullString{String: filePath, Valid: true},
				SizeBytes:      sql.NullInt64{Int64: fileSize.Bytes, Valid: true},
				YoutubeVideoID: sql.NullString{String: body.VideoId, Valid: true},
				ThumbnailPath:  sql.NullString{String: thumbnailPath, Valid: true},
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		userAudio, err := CreateUserAudio(r.Context(), apiCfg.DB, authPayload.UserId, newAudio.ID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := struct {
			*UserAudioDto
			Audio *AudioDto `json:"audio"`
		}{
			UserAudioDto: userAudioEntityToDto(userAudio),
			Audio:        audioEntityToDto(newAudio),
		}

		shared.ResCreated(w, res)
	}
}
