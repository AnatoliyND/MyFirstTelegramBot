package files

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"TelegramBot/MyFirstTgBot/lib/e"
	"TelegramBot/MyFirstTgBot/storage"
)

type Storage struct { //определяем тип в котором будет реализовываться интерфейс
	basePath string //хранит информацию о том в какой папке мы все то будем хранить
}

const defaultPerm = 0774 //параметры доступа для созданой дирректории. параметры определяются в восьмиричном порядке. 0774 - у всех пользователей будут права для чтения и записей

func New(basePath string) Storage { //передаем базовый путь при создании и возвращаем Storage
	return Storage{basePath: basePath}
}

func (s Storage) Save(_ context.Context, page *storage.Page) (err error) { //создаем метод Save
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName) //путь куда будет сохраняться наш файл. Функция filepath.Join форматизирует путь не зависимо от того где будет ипользована программа(Mac ("/"), Linux("/"), Windows("\"))

	if err := os.MkdirAll(fPath, defaultPerm); err != nil { //создаем все дирректории которые входят в корректный путь
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath) //создаем файл и передаем путь до файла
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil { //преобразует страницу в формат gob и записывает в указаный файл
		return err
	}

	return nil
}

func (s Storage) PickRandom(_ context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName) //получаем путь до дерриктории с файлами

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 { //если файлов оказалось 0
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano()) //псевдослучайное число, для выбора нового числа при каждом запуске используем время
	n := rand.Intn(len(files))       //генерируем случайное число, верхняя граница - количество файлов(ltn(files))

	file := files[n] //получаем случайный файл, с тем номером, который сгенерировали

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(_ context.Context, p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(_ context.Context, p *storage.Page) (bool, error) { //этот метод возвращает логический параметр, который будет говорить о том существует ли данная страница или нет, т.е. сохранял ли ее пользователь ранее
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exist", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName) //получаем путь до файла

	switch _, err = os.Stat(path); { //проверяем существует ли файл при помощи свитча
	case errors.Is(err, os.ErrNotExist): //проверка о несуществующем файле
		return false, nil
	case err != nil: //обрабатываем все остальные ошибки
		msg := fmt.Sprintf("can't check if file %s exist", path) //добавляем путь до файла

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page //создаем переменную, в которую файл будет интегрирован

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
