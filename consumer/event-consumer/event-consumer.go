package event_consumer

import (
	"context"
	"log"
	"time"

	"TelegramBot/MyFirstTgBot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int //размер пачки, указывает сколько событий обрабатывается за один раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error { //реализайия метода старт(Start)
	for { //вечный цикл, который постоянно ждет новые события и обрабатывает их
		gotEvents, err := c.fetcher.Fetch(context.Background(), c.batchSize) //получаем событие с помощью fetcher
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 { //если получено 0 событий, пропускаем иттерацию, но ждем 1 сек
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error { //на входе получаем событие, возвращаем ошибку
	for _, event := range events { //перебираем события
		log.Printf("got new event: %s", event.Text) //пишем лог, что получили новое событие и готовы его обработать

		if err := c.processor.Process(ctx, event); err != nil { //если с обработкой пойдет что-то не так, то переходим к следующей иттерации
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
