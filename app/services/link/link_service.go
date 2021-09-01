package link

import (
	"fmt"
	"net/http"
	"time"
	"urlshortener/app/services/shortener"

	uuid "github.com/satori/go.uuid"
)

type LinkStorer interface {
	Insert(link *Link) error
	Select(query string) ([]Link, error)
	Update(query string) error
	Exist(userID *uuid.UUID, longLink string) (bool, error)
}

type LinkService struct {
	Store         LinkStorer
	serverAddress string
}

func NewLinkService(store LinkStorer, serverAddress string) *LinkService {
	return &LinkService{Store: store, serverAddress: serverAddress}
}

func (l *LinkService) CreateLink(userID *uuid.UUID, longLink string) (string, error) {
	ok, err := l.Store.Exist(userID, longLink)
	if err != nil {
		return "", err
	}

	if ok {
		return "", fmt.Errorf("this link already exist")
	}

	link := &Link{
		ID:        uuid.NewV4(),
		OwnerID:   *userID,
		CreatedAt: time.Now(),
		LongLink:  longLink,
		ShortLink: l.createShortLink(),
	}

	err = l.Store.Insert(link)
	if err != nil {
		return "", err
	}

	return link.ShortLink, nil
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

func (l *LinkService) createShortLink() string {
	return l.serverAddress + "/" + shortener.Shorten()
}
