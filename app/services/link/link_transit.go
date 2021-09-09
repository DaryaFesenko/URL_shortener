package link

import "github.com/google/uuid"

type LinkTransition struct {
	ID         uuid.UUID `json:"id"`
	LinkID     uuid.UUID `json:"linkId"`
	UsedUserID string    `json:"usedUserId"` // заменить на другой уникальный параметр string
	UsedCount  int       `json:"usedCount"`
}
