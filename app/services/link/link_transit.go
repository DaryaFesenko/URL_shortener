package link

import (
	"time"

	"github.com/google/uuid"
)

type LinkTransition struct {
	ID        uuid.UUID `json:"id"`
	Date      time.Time `json:"date"`
	LinkID    uuid.UUID `json:"linkId"`
	IP        string    `json:"ip"`
	UsedCount int       `json:"usedCount"`
}
