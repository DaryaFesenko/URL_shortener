package linktransition

import (
	"database/sql"
	"net/http"
)

type LinkTransitService struct {
	Store *LinkTransitionStore
}

func NewLinkTransitService(db *sql.DB) *LinkTransitService {
	store := NewLinkTransitStorer(db)

	return &LinkTransitService{Store: store}
}

func (l *LinkTransitService) StatisticLink(w http.ResponseWriter, r *http.Request) {
	//получить юзера
	//получить id ссылки
	//проверить что ссылка пренадлежит юзеру

	//собрать статистику
	// сформировать ответ
}
