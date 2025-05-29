package audio

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
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

		form, err := ValidateUploadUserLocalMusicDTO(w, r)
		if err != nil {
			shared.ResTryHttpError(w, err)
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

func handleDeleteUnusedAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		unusedAudios, err := GetUnusedAudios(
			ctx,
			apiCfg.DB,
		)

		if err != nil {
			log.Println("Error deleting unused audios: ", err)
			return
		}

		for _, unusedAudio := range unusedAudios {
			if err := shared.RunNoResultDBTransaction(
				ctx,
				apiCfg.ResourceConfig,
				func(tx *database.Queries) error {
					if err := DeleteAudioByID(
						context.Background(),
						apiCfg.DB,
						unusedAudio.ID,
					); err != nil {
						log.Println("Error deleting audio: ", err)
						return err
					}

					if err := os.Remove(unusedAudio.Path.String); err != nil {
						log.Println("Error removing unused audio file: ", err)
						return err
					}

					if unusedAudio.ThumbnailPath.Valid {
						if err := os.Remove(unusedAudio.ThumbnailPath.String); err != nil {
							log.Println("Error removing unused audio thumbnail file: ", err)
							return err
						}
					}

					return nil
				},
			); err != nil {
				log.Println("Error deleting unused audio: ", err)
			}
		}

		res := deleteUnusedAudioResultDTO{
			DeletedCount: len(unusedAudios),
			AudioNames: shared.Map(unusedAudios, func(audio database.Audio) string {
				return audio.Title.String
			}),
		}

		shared.ResOK(w, res)
	}
}
