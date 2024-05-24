package audiolike

import "github.com/kume1a/sonifybackend/internal/database"

func AudioLikeEntityToDTO(e *database.AudioLike) *AudioLikeDTO {
	if e == nil {
		return nil
	}

	return &AudioLikeDTO{
		UserID:  e.UserID,
		AudioID: e.AudioID,
	}
}

func AudioLikeEntityListToDTOList(e []database.AudioLike) []*AudioLikeDTO {
	res := make([]*AudioLikeDTO, len(e))
	for i, v := range e {
		res[i] = AudioLikeEntityToDTO(&v)
	}

	return res
}
