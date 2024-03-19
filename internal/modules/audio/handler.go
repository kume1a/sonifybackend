package audio

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadYoutubeAudio(apiCfg *shared.ApiConfg) http.HandlerFunc {
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
		if _, err := GetUserAudioByYoutubeVideoId(apiCfg.DB, r.Context(), authPayload.UserId, body.VideoId); err == nil {
			shared.ResConflict(w, shared.ErrAudioAlreadyExists)
			return
		}

		videoTitle, err := youtube.GetYoutubeVideoTitle(body.VideoId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		audioDurationInSeconds, err := youtube.GetYoutubeAudioDurationInSeconds(body.VideoId)
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
			apiCfg.DB,
			r.Context(),
			sql.NullString{String: videoTitle, Valid: true},
			sql.NullString{},
			sql.NullInt32{Int32: int32(audioDurationInSeconds), Valid: true},
			filePath,
			authPayload.UserId,
			sql.NullInt64{Int64: fileSize.Bytes, Valid: true},
			sql.NullString{String: body.VideoId, Valid: true},
			sql.NullString{String: thumbnailPath, Valid: true},
		)
		if err != nil {
			log.Println("Error creating audio: ", err)
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := audioEntityToDto(newAudio)
		shared.ResCreated(w, res)
	}
}
