package hiddenuseraudio

import "github.com/kume1a/sonifybackend/internal/database"

func HiddenUserAudioEntityToDTO(e *database.HiddenUserAudio) *HiddenUserAudioDTO {
	if e == nil {
		return nil
	}

	return &HiddenUserAudioDTO{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,

		UserID:  e.UserID,
		AudioID: e.AudioID,
	}
}

func HiddenUserAudioEntityListToDTOList(e []database.HiddenUserAudio) []*HiddenUserAudioDTO {
	res := make([]*HiddenUserAudioDTO, len(e))
	for i, v := range e {
		res[i] = HiddenUserAudioEntityToDTO(&v)
	}

	return res
}
