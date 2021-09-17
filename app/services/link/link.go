package link

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	ShortLink string `json:"shortLink"`
	LongLink  string `json:"longLink"`

	OwnerID uuid.UUID `json:"ownerId"`
}
