package link

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Link struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	ShortLink string `json:"shortLink"`
	LongLink  string `json:"longLink"`

	OwnerID uuid.UUID `json:"ownerId"`
}
