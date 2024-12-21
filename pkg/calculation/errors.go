package calculation

import "errors"

// Набор ошибок калькулятора
var (
	ErrEmptyInput				 = errors.New("пустое выражение")
	ErrUnopenedParen		 = errors.New("неоткрытая скобка")
	ErrUnclosedParen		 = errors.New("незакрытая скобка")
	ErrTooManyOperats	 	 = errors.New("слишком много операторов")
	ErrDivisionByZero    = errors.New("деление на ноль")
	ErrParsing					 = errors.New("не удалось преобразовать в число, возможно не все символы цифры или операторы")

	KnownErrors					 = []error{ErrEmptyInput, ErrUnclosedParen, ErrUnopenedParen, ErrTooManyOperats, ErrDivisionByZero, ErrParsing}
)

