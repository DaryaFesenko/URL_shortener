package link_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
	"urlshortener/app/services/link"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testLink1 = "http://testlink1"
)

var _ link.LinkStorer = &Store{}

type Store struct {
	links []*link.Link
}

func InitStore() Store {
	s := Store{}
	return s
}

func (s *Store) Insert(link *link.Link) error {
	s.links = append(s.links, link)
	return nil
}

func (s *Store) ExistLongLink(userID *uuid.UUID, longLink string) (bool, error) {
	for i := range s.links {
		link := s.links[i]
		if link.OwnerID == *userID && link.LongLink == longLink {
			return true, nil
		}
	}
	return false, nil
}

func (s *Store) GetLinks(ownerID uuid.UUID) ([]link.Link, error) {
	return []link.Link{}, nil
}

func (s *Store) GetLink(linkID uuid.UUID) (*link.Link, error) {
	for i := range s.links {
		link := s.links[i]
		if link.ID == linkID {
			return link, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (s *Store) GetLongLink(shortLink string) (string, error) {
	for i := range s.links {
		link := s.links[i]
		if link.ShortLink == shortLink {
			return link.LongLink, nil
		}
	}

	return "", fmt.Errorf("can't get long link with short link: %s", shortLink)
}

func (s *Store) GetLinkIDByShortLink(shortLink string) (*uuid.UUID, error) {
	for i := range s.links {
		link := s.links[i]
		if link.ShortLink == shortLink {
			return &link.ID, nil
		}
	}

	return nil, fmt.Errorf("can't get link id with short link: %s", shortLink)
}

func (s *Store) DeleteLink(linkID uuid.UUID) error {
	for i := range s.links {
		link := s.links[i]
		if link.ID == linkID {
			s.links = append(s.links[:i], s.links[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("can't get link with id: %s", linkID)
}

func TestCreateLinkOK(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	_, err := service.CreateLink(&userId, testLink1)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(store.links), 1)
}

func TestDeleteLinkOK(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	linkId := uuid.New()
	store.Insert(&link.Link{
		ID:        linkId,
		CreatedAt: time.Now(),
		ShortLink: "",
		LongLink:  testLink1,
		OwnerID:   userId,
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkId,
		IP:        "user1",
		UsedCount: 2,
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkId,
		IP:        "user2",
		UsedCount: 1,
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    uuid.New(),
		IP:        "user3",
		UsedCount: 1,
	})

	err := service.DeleteLink(userId, linkId)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(linkTransitStore.linkTransitions), 1)
	assert.Equal(t, len(store.links), 0)
}

func TestDeleteLinkFAIL(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	linkId := uuid.New()

	err := service.DeleteLink(userId, linkId)

	assert.Equal(t, err, sql.ErrNoRows)
}

func TestGetLinkStatisticOK(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	linkId := uuid.New()

	store.Insert(&link.Link{
		ID:        linkId,
		CreatedAt: time.Now(),
		ShortLink: "",
		LongLink:  testLink1,
		OwnerID:   userId,
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkId,
		IP:        "user1",
		UsedCount: 2,
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkId,
		IP:        "user2",
		UsedCount: 1,
	})

	stat, err := service.GetLinkStatistic(&userId, linkId)
	assert.Equal(t, err, nil)
	assert.Equal(t, stat.TransitCount, 3)
	assert.Equal(t, stat.UniqueTransitCount, 2)
}

func TestGetLinkStatisticFAIL1(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	linkId := uuid.New()

	_, err := service.GetLinkStatistic(&userId, linkId)
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestGetLinkStatisticFAIL2(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userId := uuid.New()
	linkId := uuid.New()

	store.Insert(&link.Link{
		ID:        linkId,
		CreatedAt: time.Now(),
		ShortLink: "",
		LongLink:  testLink1,
		OwnerID:   uuid.New(),
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkId,
		IP:        "user1",
		UsedCount: 2,
	})

	_, err := service.GetLinkStatistic(&userId, linkId)
	assert.Equal(t, err, fmt.Errorf("link %s does not belong user %s", linkId, userId))
}

func TestGetLongLinkOK1(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userID := "user1"
	store.Insert(&link.Link{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		ShortLink: "shortLink1",
		LongLink:  testLink1,
		OwnerID:   uuid.New(),
	})

	longLink, err := service.GetLongLink("shortLink1", userID)
	assert.Equal(t, err, nil)
	assert.Equal(t, longLink, testLink1)
	assert.Equal(t, len(linkTransitStore.linkTransitions), 1)
	assert.Equal(t, linkTransitStore.linkTransitions[0].UsedCount, 1)
	assert.Equal(t, linkTransitStore.linkTransitions[0].IP, userID)
}

func TestGetLongLinkOK2(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userID := "user1"
	linkID := uuid.New()
	store.Insert(&link.Link{
		ID:        linkID,
		CreatedAt: time.Now(),
		ShortLink: "shortLink1",
		LongLink:  testLink1,
		OwnerID:   uuid.New(),
	})

	linkTransitStore.Insert(link.LinkTransition{
		ID:        uuid.New(),
		LinkID:    linkID,
		IP:        userID,
		UsedCount: 2,
	})

	longLink, err := service.GetLongLink("shortLink1", userID)
	assert.Equal(t, err, nil)
	assert.Equal(t, longLink, testLink1)
	assert.Equal(t, len(linkTransitStore.linkTransitions), 1)
	assert.Equal(t, linkTransitStore.linkTransitions[0].UsedCount, 3)
	assert.Equal(t, linkTransitStore.linkTransitions[0].IP, userID)
}

func TestGetLongLinkFAIL1(t *testing.T) {
	store := InitStore()
	linkTransitStore := InitLinkTransitStore()
	service := link.NewLinkService(&store, &linkTransitStore)

	userID := "user1"
	store.Insert(&link.Link{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		ShortLink: "shortLink1",
		LongLink:  testLink1,
		OwnerID:   uuid.New(),
	})

	_, err := service.GetLongLink("shortLink2", userID)
	assert.Equal(t, err, fmt.Errorf("can't get long link with short link: %s", "shortLink2"))
}
