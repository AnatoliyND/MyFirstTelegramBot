package config

import (
	"flag"
	"log"
)

type Config struct {
	TgBotHost string
}

func MustLoad() Config { //функция для получения токена, must(обычно применяется для запуска программы и парсенга конфигов) - предупреждает, что функция написана без обработки ошибок(стараться так не делать)
	// во время запуска программы флаг указывается в виде записи bot -tg-bot-token 'my token
	tgBotToken := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	) //первый аргумент указываем имя флага, вторым аргументом указываем значение флага по умолчанию, последним аргумент подсказка к данному флагу. В токене будет лежать не значение, а ссылка на значение, которое выпадает во время вызова функции Parse

	flag.Parse()

	if *tgBotToken == "" { //проверяем что в токене что-то лежит
		log.Fatal("token is not specified")
	}
	return Config{
		TgBotHost: *tgBotToken,
	}
}
