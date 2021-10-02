package link

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"urlshortener/app/services/shortener"

	"github.com/google/uuid"
)

type LinkStatistic struct {
	UniqueTransitCount int    `json:"uniqueCount"`
	TransitCount       int    `json:"allCount"`
	LongLink           string `json:"longLink"`
	ShortLink          string `json:"shortLink"`
}

type LinkService struct {
	store               LinkStorer
	linkTransitionStore LinkTransitStorer
}

func NewLinkService(store LinkStorer, linkTransitStore LinkTransitStorer) *LinkService {
	return &LinkService{store: store, linkTransitionStore: linkTransitStore}
}

func (l *LinkService) GetLinks(userID uuid.UUID) ([]Link, error) {
	return l.store.GetLinks(userID)
}

func (l *LinkService) CreateLink(userID *uuid.UUID, longLink string) (string, error) {
	if !strings.Contains(longLink, "http") {
		return "", fmt.Errorf("string is not link")
	}

	ok, err := l.store.ExistLongLink(userID, longLink)
	if err != nil {
		return "", err
	}

	if ok {
		return "", fmt.Errorf("this link already exist")
	}

	link := &Link{
		ID:        uuid.New(),
		OwnerID:   *userID,
		CreatedAt: time.Now(),
		LongLink:  longLink,
		ShortLink: shortener.Shorten(),
	}

	err = l.store.Insert(link)
	if err != nil {
		return "", err
	}

	return link.ShortLink, nil
}

func (l *LinkService) DeleteLink(userID, linkID uuid.UUID) error {
	link, err := l.store.GetLink(linkID)
	if err != nil {
		return err
	}

	if link.OwnerID != userID {
		return fmt.Errorf("link %s does not belong user %s", linkID, userID)
	}

	err = l.linkTransitionStore.DeleteLinkTransit(linkID)
	if err != nil {
		return err
	}

	err = l.store.DeleteLink(linkID)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkService) GetLinkStatistic(userID *uuid.UUID, linkID uuid.UUID) (LinkStatistic, error) {
	link, err := l.store.GetLink(linkID)
	if err != nil {
		return LinkStatistic{}, err
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return LinkStatistic{}, fmt.Errorf("can't found link with id: %s", linkID)
		}

		return LinkStatistic{}, err
	}

	if link.OwnerID != *userID {
		return LinkStatistic{}, fmt.Errorf("link %s does not belong user %s", linkID, userID)
	}

	linkTransitions, err := l.linkTransitionStore.StatisticLink(linkID)
	if err != nil {
		return LinkStatistic{}, err
	}

	var count int

	for i := range linkTransitions {
		count = count + linkTransitions[i].UsedCount
	}

	return LinkStatistic{
		UniqueTransitCount: len(linkTransitions),
		TransitCount:       count,
		LongLink:           link.LongLink,
		ShortLink:          link.ShortLink,
	}, nil
}

func (l *LinkService) GetLongLink(shortLink string, userID string) (string, error) {
	longLink, err := l.store.GetLongLink(shortLink)

	if err != nil {
		return "", err
	}

	linkID, err := l.store.GetLinkIDByShortLink(shortLink)
	if err != nil {
		return "", err
	}

	linkTransit, err := l.linkTransitionStore.GetTransit(userID, *linkID)
	if err != nil {
		if err == sql.ErrNoRows {
			errInsert := l.linkTransitionStore.Insert(LinkTransition{
				ID:         uuid.New(),
				LinkID:     *linkID,
				UsedUserID: userID,
				UsedCount:  1,
			})

			if errInsert != nil {
				return "", errInsert
			}

			return longLink, nil
		}

		return "", err
	}

	linkTransit.UsedCount++
	err = l.linkTransitionStore.UpdateTransitCount(linkTransit.ID, linkTransit.UsedCount)
	if err != nil {
		return "", err
	}

	return longLink, nil
}
