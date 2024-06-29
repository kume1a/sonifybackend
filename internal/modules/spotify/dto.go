package spotify

import "github.com/google/uuid"

type spotifyAccessTokenDTO struct {
	SpotifyAccessToken string `json:"spotifyAccessToken" valid:"required"`
}

type refreshSpotifyTokenDTO struct {
	SpotifyRefreshToken string `json:"spotifyRefreshToken" valid:"required"`
}

type downloadSpotifyPlaylistDTO struct {
	SpotifyAccessToken string `json:"spotifyAccessToken" valid:"required"`
	SpotifyPlaylistID  string `json:"spotifyPlaylistId" valid:"required"`
}

type authorizeSpotifyDTO struct {
	Code string `json:"code" valid:"required"`
}

type spotifyAuthCodeTokenPayloadDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type spotifyRefreshTokenPayloadDTO struct {
	AccessToken string `json:"accessToken"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type searchSpotifyQueryDTO struct {
	Keyword            []string `json:"keyword" valid:"required"`
	SpotifyAccessToken []string `json:"spotifyAccessToken" valid:"required"`
}

type searchSpotifyResDTO struct {
	Playlists []searchSpotifyResPlaylistDTO `json:"playlists"`
}

type searchSpotifyResPlaylistDTO struct {
	Name       string     `json:"name"`
	ImageUrl   string     `json:"imageUrl"`
	SpotifyID  string     `json:"spotifyId"`
	PlaylistID *uuid.UUID `json:"playlistId"`
}

// --------- Spotify API DTOs ---------

type spotifyClientCredsTokenDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type spotifyAuthCodeTokenDTO struct {
	spotifyClientCredsTokenDTO
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type spotifyRefreshTokenDTO struct {
	spotifyClientCredsTokenDTO
	Scope string `json:"scope"`
}

type spotifyPlaylistDTO struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalURLs  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		DisplayName  string `json:"display_name"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"owner"`
	PrimaryColor string `json:"primary_color"`
	Public       bool   `json:"public"`
	SnapshotID   string `json:"snapshot_id"`
	Tracks       struct {
		Href  string `json:"href"`
		Items []struct {
			AddedAt string `json:"added_at"`
			AddedBy struct {
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"added_by"`
			IsLocal      bool        `json:"is_local"`
			PrimaryColor interface{} `json:"primary_color"`
			Track        struct {
				Album struct {
					AlbumType string `json:"album_type"`
					Artists   []struct {
						ExternalURLs struct {
							Spotify string `json:"spotify"`
						} `json:"external_urls"`
						Href string `json:"href"`
						ID   string `json:"id"`
						Name string `json:"name"`
						Type string `json:"type"`
						URI  string `json:"uri"`
					} `json:"artists"`
					AvailableMarkets []string `json:"available_markets"`
					ExternalURLs     struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href   string `json:"href"`
					ID     string `json:"id"`
					Images []struct {
						Height int    `json:"height"`
						URL    string `json:"url"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name                 string `json:"name"`
					ReleaseDate          string `json:"release_date"`
					ReleaseDatePrecision string `json:"release_date_precision"`
					TotalTracks          int    `json:"total_tracks"`
					Type                 string `json:"type"`
					URI                  string `json:"uri"`
				} `json:"album"`
				Artists []struct {
					ExternalURLs struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				DiscNumber       int      `json:"disc_number"`
				DurationMS       int      `json:"duration_ms"`
				Episode          bool     `json:"episode"`
				Explicit         bool     `json:"explicit"`
				ExternalIDs      struct {
					ISRC string `json:"isrc"`
				} `json:"external_ids"`
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href        string `json:"href"`
				ID          string `json:"id"`
				IsLocal     bool   `json:"is_local"`
				Name        string `json:"name"`
				Popularity  int    `json:"popularity"`
				PreviewURL  string `json:"preview_url"`
				Track       bool   `json:"track"`
				TrackNumber int    `json:"track_number"`
				Type        string `json:"type"`
				URI         string `json:"uri"`
			} `json:"track"`
			VideoThumbnail struct {
				URL string `json:"url"`
			} `json:"video_thumbnail"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type spotifyGetPlaylistsDTO struct {
	Href     string               `json:"href"`
	Limit    int                  `json:"limit"`
	Next     string               `json:"next"`
	Offset   int                  `json:"offset"`
	Previous string               `json:"previous"`
	Total    int                  `json:"total"`
	Items    []spotifyPlaylistDTO `json:"items"`
}

type spotifyPlaylistItemDTO struct {
	Track struct {
		Album struct {
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"album"`
		Artists []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"artists"`
		DurationMS int    `json:"duration_ms"`
		Name       string `json:"name"`
		PreviewURL string `json:"preview_url"`
		ID         string `json:"id"`
	} `json:"track"`
	AddedBy struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"added_by"`
	AddedAt string `json:"added_at"`
}

type spotifyPlaylistItemsDTO struct {
	Limit    int                      `json:"limit"`
	Items    []spotifyPlaylistItemDTO `json:"items"`
	Next     string                   `json:"next"`
	Total    int                      `json:"total"`
	Previous string                   `json:"previous"`
}
type spotifySearchDTO struct {
	Playlists struct {
		Href     string                         `json:"href"`
		Limit    int                            `json:"limit"`
		Next     string                         `json:"next"`
		Offset   int                            `json:"offset"`
		Previous string                         `json:"previous"`
		Total    int                            `json:"total"`
		Items    []spotifySearchPlaylistItemDTO `json:"items"`
	} `json:"playlists"`
}

type spotifySearchPlaylistItemDTO struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalURLs  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"followers"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Public     bool   `json:"public"`
	SnapshotID string `json:"snapshot_id"`
	Tracks     struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}
