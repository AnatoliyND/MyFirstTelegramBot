package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"TelegramBot/MyFirstTgBot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error                        //принимает ссылку на страницу и передает ошибку
	PickRandom(ctx context.Context, userName string) (*Page, error) //принимает имя пользователя и возвращает страницу
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error) //проверяет существует ли та или иная страница
}

var ErrNoSavedPages = errors.New("no saved page")

type Page struct { //основной тип данных с которым будет работать Storage. Будет принимать страницу на которую ведет ссылка которую будем передавать боту
	URL      string
	UserName string //имя пользователя который скинул ссылку(для понимания кому ее отдавать)
}

func (p Page) Hash() (string, error) { //создаем метод для генерации названия файла(в котором будут храниться ссылки). Возвращает текстовое представление хэша и ошибку
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
