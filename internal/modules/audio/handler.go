package audio

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleUploadUserLocalMusic(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		form, httpErr := ValidateUploadUserLocalMusicDTO(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		audioExists, err := DoesAudioExistByLocalId(r.Context(), apiCfg.DB, authPayload.UserID, form.LocalID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if audioExists {
			shared.DeleteFiles([]string{form.AudioPath, form.ThumbnailPath})

			shared.ResConflict(w, shared.ErrAudioAlreadyExists)
			return
		}

		userAudioWithAudio, httpErr := WriteUserImportedLocalMusic(
			WriteUserImportedLocalMusicParams{
				ResourceConfig:     apiCfg.ResourceConfig,
				Context:            r.Context(),
				UserID:             authPayload.UserID,
				AudioTitle:         form.Title,
				AudioAuthor:        form.Author,
				AudioPath:          form.AudioPath,
				AudioThumbnailPath: form.ThumbnailPath,
				AudioDurationMs:    form.DurationMs,
				AudioLocalId:       form.LocalID,
			},
		)

		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		res := sharedmodule.UserAudioWithRelDTO{
			UserAudioDTO: sharedmodule.UserAudioEntityToDTO(userAudioWithAudio.UserAudio),
			Audio:        sharedmodule.AudioEntityToDto(userAudioWithAudio.Audio),
		}

		shared.ResCreated(w, res)
	}
}

func handleWriteInitialAudioRelCount(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := WriteInitialAudioRelCount(r.Context(), apiCfg.ResourceConfig); err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		shared.ResOK(w, nil)
	}
}
