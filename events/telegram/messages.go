package telegram

//Команда ниже содержит в себе краткую справку
const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.
❗❗❗ Caution! After that, this page will be removed from your list! 🚮` //предупреждение о том что ссылка после команды /rnd будет удалена

//Команда ниже содержит в себе краткую справку и небольшое приветствие. Потребуется для старта общения с пользователем

const msgHello = "Hi there! 👋\n\n" + msgHelp

//Пишем группу сообщений которыми бот будет комментировать наши сообщения

const (
	msgUnknownCommand = "Unknown command🤨"                              //Сообщение на случай неизвестной команды
	msgNoSavedPages   = "You have not saved pages🍩"                     //это сообщение на случай когда пользователь запрашивает ссылку, но у него уже не осталось непрочтенных ссылок
	msgSaved          = "✅Saved!👌"                                      //собщение, когда пользователь скидывает ссылку и бот успешно ее сохраняет
	msgAlreadyExists  = "You have already have this page in your list🤖" //сообщение на тот случай, когда пользователь пытается сохранить ссылку, которую сохранял ранее и она все еще хранится в списке
)
