package link

import (
	"database/sql"
	"net/http"
)

type LinkService struct {
	Store *LinkStore
}

func NewLinkService(db *sql.DB) *LinkService {
	store := NewLinkStorer(db)

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
