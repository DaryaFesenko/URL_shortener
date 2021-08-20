package link

import (
	"net/http"
)

type LinkStorer interface {
	Insert(link Link) error
	Select(query string) ([]Link, error)
	Update(query string) error
}

type LinkService struct {
	Store LinkStorer
}

func NewLinkService(store LinkStorer) *LinkService {
	return &LinkService{Store: store}
}

func (l *LinkService) CreateLink(w http.ResponseWriter, r *http.Request) {
	// найти юзера из авторизации

	// получить длинную ссылку из запроса

	//проверить что такой ссылки нет

	//создать ссылку

	//вернуть короткую ссылку
}

func (l *LinkService) DeleteLink(w http.ResponseWriter, r *http.Request) {
	// найти юзера из авторизации

	// получить id ссылки из запроса

	//проверить что ссылка принадлежит юзеру

	//удалить переходы по ссылке

	//удалить ссылку
}

func (l *LinkService) GetLongLink(w http.ResponseWriter, r *http.Request) {
	//получить какой нибудь уникальный идентификатор типа адреса

	//получить короткую ссылку из запроса

	//создать запись о переходе если ее нет
	//если есть, плюсануть переход

	// вернуть длинную ссылку
}

func (l *LinkService) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	// получить юзера

	//получить все ссылки юзера

	//сформировать ответ
}
