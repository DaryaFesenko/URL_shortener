package linktransit

import (
	"net/http"
)

type LinkTransitStorer interface {
	Insert(lt LinkTransition) error
	Select(query string) ([]LinkTransition, error)
	Update(query string) error
}

type LinkTransitService struct {
	Store LinkTransitStorer
}

func NewLinkTransitService(store LinkTransitStorer) *LinkTransitService {
	return &LinkTransitService{Store: store}
}

func (l *LinkTransitService) StatisticLink(w http.ResponseWriter, r *http.Request) {
	//получить юзера
	//получить id ссылки
	//проверить что ссылка пренадлежит юзеру

	//собрать статистику
	// сформировать ответ
}
