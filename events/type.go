package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, e Event) error //метод Process принимает событие (Event) и возвращает ошибку
}

type Type int

const ( //создаем список событий которые будем использовать
	Unknown Type = iota //неизвестный тип для обработки случаев, когда не смогли определить тип события. iota используется для определения групп констант, первой константе она присваевает 0 и дальше идет по порядку
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
