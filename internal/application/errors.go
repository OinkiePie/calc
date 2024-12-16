package application

import "errors"

// Набор ошибок сервера
var (
	ErrOnlyPostAllowed	 = errors.New("только запросы типа POST разрешены")
	ErrFailedToUnmarshal = errors.New("не удалось демаршалировать запрос")
	ErrEmptyRequest			 = errors.New("получен пустой запрос")
	ErrInvalidChars		 	 = errors.New("некорректные символы")
)