package playlist

import (
	"database/sql"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ValidateCreatePlaylistDto(w, r)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		playlist, err := CreatePlaylist(r.Context(), apiCfg.DB, database.CreatePlaylistParams{
			Name:              body.Name,
			ThumbnailPath:     sql.NullString{String: body.ThumbnailPath, Valid: true},
			AudioImportStatus: database.ProcessStatusCOMPLETED,
			AudioCount:        0,
			TotalAudioCount:   0,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := PlaylistEntityToDto(*playlist)

		shared.ResCreated(w, dto)
	}
}

func handleGetPlaylists(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := shared.ValidateRequestQuery[*shared.LastCreatedAtPageParamsDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetPlaylists(r.Context(), apiCfg.DB, database.GetPlaylistsParams{
			CreatedAt: query.LastCreatedAt,
			Limit:     query.Limit,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, PlaylistEntityToDto)

		shared.ResOK(w, dtos)
	}
}

func handleCreatePlaylistAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*createPlaylistAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlistAudio, err := playlistaudio.CreatePlaylistAudio(r.Context(), apiCfg.DB, database.CreatePlaylistAudioParams{
			PlaylistID: body.PlaylistID,
			AudioID:    body.AudioID,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := playlistAudioEntityToDto(playlistAudio)

		shared.ResCreated(w, dto)
	}
}

func handleGetPlaylistWithAudios(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPaylad, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlistIDDTO, err := ValidateGetPlaylistByIDVars(r)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		playlist, audios, err := GetPlaylistWithAudios(r.Context(), apiCfg.DB, playlistIDDTO.PlaylistID, authPaylad.UserID)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		dto := struct {
			PlaylistDTO
			Audios []*audio.AudioDTO `json:"audios"`
		}{
			PlaylistDTO: PlaylistEntityToDto(*playlist),
			Audios:      shared.Map(audios, audio.AudioWithAudioLikeToAudioDTO),
		}

		shared.ResOK(w, dto)
	}
}
