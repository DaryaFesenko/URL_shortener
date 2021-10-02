package link

import "github.com/google/uuid"

type LinkStorer interface {
	Insert(link *Link) error
	ExistLongLink(userID *uuid.UUID, longLink string) (bool, error)
	GetLink(linkID uuid.UUID) (*Link, error)
	GetLongLink(shortLink string) (string, error)
	GetLinkIDByShortLink(shortLink string) (*uuid.UUID, error)
	DeleteLink(linkID uuid.UUID) error
	GetLinks(ownerID uuid.UUID) ([]Link, error)
}

type LinkTransitStorer interface {
	Insert(lt LinkTransition) error
	UpdateTransitCount(id uuid.UUID, usedCount int) error
	StatisticLink(linkID uuid.UUID) ([]LinkTransition, error)
	GetTransit(ip string, linkID uuid.UUID) (LinkTransition, error)
	DeleteLinkTransit(linkID uuid.UUID) error
}
