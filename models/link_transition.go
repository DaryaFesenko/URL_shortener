package models

import uuid "github.com/satori/go.uuid"

type LinkTransition struct {
	ID         uuid.UUID `json:"id"`
	LinkID     uuid.UUID `json:"linkId"`
	UsedUserID uuid.UUID `json:"usedUserId"` // заменить на другой уникальный параметр
	UsedCount  int       `json:"usedCount"`
}
