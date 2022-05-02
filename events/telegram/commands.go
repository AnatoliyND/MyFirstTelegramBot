package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"TelegramBot/MyFirstTgBot/lib/e"
	"TelegramBot/MyFirstTgBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error { //метод типа процессор. Это что-то вроде API роутера. Смотрим на текст сообщения и по его формату и содержинию будем понимать какая это команда. Сюда передаем текст, чат ID и юзернейм, чтобы понимать что с этим делать и куда отправлять сообщение
	text = strings.TrimSpace(text) //для начала удаляем с текста лишние пробелы

	log.Printf("get new command '%s' from '%s'", text, username) //получаем новую команду, ее содержимое и кто автор сообщения

	if isAdCmd(text) { //проверяем является ли сообщение ссылкой
		return p.savePage(ctx, chatID, text, username)
	}

	switch text { //расставляем действия
	case RndCmd:
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{ //подготавливаем страницу которую собираемся сохранить
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil { //пытаемся сохранить страницу
		return err
	}

	if err := p.tg.SendMessage(ctx, chatID, msgSaved); err != nil { //если страница успешно сохранилась сообщаем об этом пользователю
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) { //функция отправляет пользователю случайную статью. на входе принимаем chatID и пользователя чьи ссылки мы перебираем
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, username) //ищем случайную статью
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) { //если ошибка всеже есть, то возвращаем пользователю сообщение, что он ничего не сохранил
		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil { //если боту надолось что-то найти, то он отправляет ссылку пользователю
		return err
	}

	return p.storage.Remove(ctx, page) //если удалось найти и отправить ссылку, то ее нужно удалить
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error { //отправка справки
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error { //отправка команды приветствия
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAdCmd(text string) bool { //функция проверяет является ли текст ссылкой
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text) //распарсим текущай текст считая его ссылкой

	return err == nil && u.Host != "" //текст мы будем считать ссылкой если ошибка нулевая и при этом указан хост(*недостаток данного подхода в том что ссылки типа ya.ru не будут считаться ссылками, т.е. те у которых не указан протакол, должен всегда присутствовать префикс (http(s)://) Пример: http://ya.ru либо https://ya.ru)
}
