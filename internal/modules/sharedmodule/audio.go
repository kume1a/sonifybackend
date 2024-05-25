package sharedmodule

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func ValidateAudioExistsByID(
	ctx context.Context,
	db *database.Queries,
	audioID uuid.UUID,
) error {
	count, err := db.CountAudioByID(ctx, audioID)

	if err != nil {
		log.Println("Error validating audio exists by ID: ", err)
		return shared.InternalServerErrorDef()
	}

	if count == 0 {
		return shared.NotFound(shared.ErrAudioNotFound)
	}

	return nil
}
