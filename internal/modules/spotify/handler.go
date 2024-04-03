package spotify

import (
	"log"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadPlaylist(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*downloadSpotifyPlaylistDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		spotifyPlaylistDTO, err := GetSpotifyPlaylist(body.SpotifyAccessToken, body.PlaylistID)
		if err != nil {
			shared.ResInternalServerError(w, shared.ErrFailedToGetSpotifyPlaylist)
			return
		}

		spotifyPlaylist := spotifyPlaylistDtoToModel(spotifyPlaylistDTO)

		trackDownloadMetas := []*downloadSpotifyTrackMetaDTO{}
		for _, track := range spotifyPlaylist.Tracks[:shared.Min(10, len(spotifyPlaylist.Tracks)-1)] {
			log.Println("Getting download meta for track: ", track.ID)

			downloadMeta, err := GetSpotifyAudioDownloadMeta(track.ID)
			if err != nil || !downloadMeta.Success {
				log.Println("failed to get download meta for track: ", track.ID)
				break
			}

			trackDownloadMetas = append(trackDownloadMetas, downloadMeta)
		}

		for _, downloadMeta := range trackDownloadMetas {
			log.Println("Downloading track: ", downloadMeta.Metadata.Title, " from: ", downloadMeta.Link)

			fileLocation, err := shared.NewPublicFileLocation(shared.PublicFileLocationArgs{
				Extension: ".mp3",
				Dir:       shared.DirSpotifyAudios,
			})
			if err != nil {
				log.Println("error creating file location: ", err)
				break
			}

			if err := shared.DownloadFile(fileLocation, downloadMeta.Link); err != nil {
				log.Println("error downloading file: ", err)
				break
			}
		}
	}
}

func handleAuthorizeSpotify(w http.ResponseWriter, r *http.Request) {
	body, err := shared.ValidateRequestBody[*authorizeSpotifyDTO](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	tokenPayload, err := GetAuthorizationCodeSpotifyTokenPayload(body.Code)
	if err != nil {
		shared.ResInternalServerError(w, shared.ErrFailedToGetSpotifyAccessToken)
		return
	}

	dto := spotifyTokenPayloadDTO{
		AccessToken:  tokenPayload.AccessToken,
		RefreshToken: tokenPayload.RefreshToken,
		Scope:        tokenPayload.Scope,
		ExpiresIn:    tokenPayload.ExpiresIn,
		TokenType:    tokenPayload.TokenType,
	}

	shared.ResOK(w, dto)
}

func handleImportSpotifyUserPlaylists(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		body, err := shared.ValidateRequestQuery[*spotifyAccessTokenDTO](r)
		if err != nil {
			log.Println("error validating request query: ", err)
			shared.ResBadRequest(w, err.Error())
			return
		}

		if err := downloadSpotifyPlaylist(apiCfg, r.Context(), authPayload.UserId, body.SpotifyAccessToken[0]); err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, struct{}{})
	}
}
