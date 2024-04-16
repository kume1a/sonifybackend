package audio

import (
	"context"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func BulkWriteAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	params []database.CreateAudioParams,
) ([]database.Audio, error) {
	return shared.RunDbTransaction(
		ctx,
		apiCfg,
		func(tx *database.Queries) ([]database.Audio, error) {
			audios := make([]database.Audio, 0, len(params))

			for _, param := range params {
				audio, err := CreateAudio(ctx, tx, param)
				if err != nil {
					return nil, err
				}

				audios = append(audios, *audio)
			}

			return audios, nil
		},
	)
}
