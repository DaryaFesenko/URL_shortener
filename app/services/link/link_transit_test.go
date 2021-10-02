package link_test

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/link"

	"github.com/google/uuid"
)

var _ link.LinkTransitStorer = &LinkTransitStore{}

type LinkTransitStore struct {
	linkTransitions []*link.LinkTransition
}

func InitLinkTransitStore() LinkTransitStore {
	s := LinkTransitStore{}
	return s
}

func (l *LinkTransitStore) Insert(lt link.LinkTransition) error {
	l.linkTransitions = append(l.linkTransitions, &lt)
	return nil
}

func (l *LinkTransitStore) UpdateTransitCount(id uuid.UUID, usedCount int) error {
	for i := range l.linkTransitions {
		link_transition := l.linkTransitions[i]
		if link_transition.ID == id {
			link_transition.UsedCount = usedCount
			return nil
		}
	}

	return fmt.Errorf("cant't update link transition with id: %s", id)
}

func (l *LinkTransitStore) StatisticLink(linkID uuid.UUID) ([]link.LinkTransition, error) {
	statistic := make([]link.LinkTransition, 0)

	for i := range l.linkTransitions {
		linkTransition := l.linkTransitions[i]

		if linkTransition.LinkID == linkID {
			statistic = append(statistic, *linkTransition)
		}
	}

	return statistic, nil
}

func (l *LinkTransitStore) GetTransit(ip string, linkID uuid.UUID) (link.LinkTransition, error) {
	for i := range l.linkTransitions {
		linkTransition := l.linkTransitions[i]

		if linkTransition.LinkID == linkID && linkTransition.IP == ip {
			return *linkTransition, nil
		}
	}

	return link.LinkTransition{}, sql.ErrNoRows
}

func (l *LinkTransitStore) DeleteLinkTransit(linkID uuid.UUID) error {
	newList := make([]*link.LinkTransition, 0)
	for i := range l.linkTransitions {
		linkTransition := l.linkTransitions[i]

		if linkTransition.LinkID != linkID {
			newList = append(newList, linkTransition)
		}
	}

	l.linkTransitions = newList

	return nil
}
