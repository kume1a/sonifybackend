package audio

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadYoutubeAudio(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[downloadYoutubeAudioDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		audioUrl, err := youtube.GetYoutubeAudioUrl(body.VideoId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
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

		log.Println("got audioUrl: " + audioUrl + " videoTitle: " + videoTitle + " audioDurationInSeconds: " + string(audioDurationInSeconds))

		fileLocation, err := shared.NewPublicFileLocation(".mp3")
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if err := shared.DownloadFile(fileLocation, audioUrl); err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		log.Println("downloaded file to: ", fileLocation)

		audio, err := CreateAudio(
			apiCfg.DB,
			r.Context(),
			sql.NullString{String: videoTitle, Valid: true},
			sql.NullString{},
			sql.NullInt32{Int32: int32(audioDurationInSeconds), Valid: true},
			fileLocation,
			uuid.New(),
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		log.Println("audio created: ", audio)

		shared.ResCreated(w, audio)
	}
}
