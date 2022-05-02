package telegram

import (
	"context"
	"errors"

	"TelegramBot/MyFirstTgBot/clients/telegram"
	"TelegramBot/MyFirstTgBot/events"
	"TelegramBot/MyFirstTgBot/lib/e"
	"TelegramBot/MyFirstTgBot/storage"
)

type Processor struct { //тип данных, который реализовывает два интерфейса
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor { //функция создает новый Processor
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) { //uodates это понятия телеграма и они относятся только к нему(в др месенджере возможно данного термина не быть). Event это более общая сущность, в нее мы можем преобразовывать все что получаем от других месенджеров, в каком бы формате они бы не предоставляли нам информацию
	updates, err := p.tg.Updates(ctx, p.offset, limit) //с помощью клиента получаем все апдейты
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 { //возвращаем нулевой результат если ивентов не нашли
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates)) //аллоцируем память под результат

	for _, u := range updates { //перебираем все апдейты и преобразовываем их в евенты
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 //берем последний апдейт, смотри его ID и добавляем к нему 1 и тогда при следующем запросе получим только те апдейты у которых ID больше чем у последнего из уже полученых

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error { //этот метод выполняет различные действия в зависимости от типа евента
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event { //функция для преобразования апдейтов в евенты
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
