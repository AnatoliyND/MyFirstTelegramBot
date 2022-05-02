package e

import "fmt"

func Wrap(msg string, err error) error { //функция для оборачивания ошибок. На входе принимает текст сообщения с подсказкой и саму ошибку, возвращет ошибку
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfErr(msg string, err error) error { //функция для оборачивания ошибок. На входе принимает текст сообщения с подсказкой и саму ошибку, возвращет ошибку
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
