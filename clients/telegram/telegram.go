package telegram

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"TelegramBot/MyFirstTgBot/lib/e"
)

type Client struct {
	host     string      //host API сервера телеграма
	basePath string      //базовый путь(это путь с которого начинаются все запросы)
	client   http.Client //тут хранится http клиент для того чтобы не создавать его для каждого запроса отдельно
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client { //функция создает клиент, в нее передает хост и токен
	return &Client{
		host:     host,
		basePath: newBasePath(token), //выносим путь в отдельную функцию
		client:   http.Client{},      //http оставляем стандартным
	}
}

func newBasePath(token string) string {
	return "bot" + token // упрощает вставку пути при повторном применении. Также в случае изменения телеграммом префикса "bot" не придется изменять большую чать кода
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (update []Update, err error) { //функция возвращает структуру, которая будет содержать все что нам нужно знать об Updates
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{} //формируем параметры запроса
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error { //метод отправки сообщений
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) { //функция для отправки запроса
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method), //путь состоит из бызовой части созданой в клиенте и метода. path.Join корректирует путь убирая или добавляя лишние слеши
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil) //формируем объект запроса. Не отправляем запрос, а только подготавливаем его
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req) //отправляем получившийся запрос
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
