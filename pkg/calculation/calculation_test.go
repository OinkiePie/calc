package calculation_test

import (
	"testing"

	"github.com/OinkiePie/calc/pkg/calculation"
)

// Проверка првильности различных вычисление
// Проверка работоспособности сервера находится в application_test.go
func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name					 string
		expression     string
		expectedResult float64
	}{
		// Базовые операции
		{
			name: "BASIC #1",
			expression: "1 + 1",
			expectedResult: 2},
		{
			name: "BASIC #2",
			expression: "10 - 5",
			expectedResult: 5,
		},
		{
			name: "BASIC #3",
			expression: "2 * 3",
			expectedResult: 6,
		},
		{
			name: "BASIC #4",
			expression: "10 / 2",
			expectedResult: 5,
		},

		// Дроби
		{
			name: "Fractions #1",
			expression: "1/2",
			expectedResult: 0.5,
		},
		{
			name: "Fractions #2",
			expression: "3/4",
			expectedResult: 0.75,
		},
		{
			name: "Fractions #3",
			expression: "10/3", expectedResult: 3.333,
		},


		// Порядок операций
		{
			name: "Priority #1",
			expression: "1 + 2 * 3",
			expectedResult: 7,
		},
		{
			name: "Priority #2",
			expression: "(1 + 2) * 3",
			expectedResult: 9,
		},
		{
			name: "Priority #3",
			expression: "10 - 2 * 3 + 1",
			expectedResult: 5,},

		// Скобки
		{
			name: "Parentheses #1",
			expression: "(1 + 2) * (3 - 1)",
			expectedResult: 6,},
		{
			name: "Parentheses #2",
			expression: "((1 + 2) * 3) - 1",
			expectedResult: 8,
		},

		// Отрицательные числа
		{
			name: "Negative #1",
			expression: "-5 + 10",
			expectedResult: 5,
		},
		{
			name: "Negative #2",
			expression: "-1 - - 2",
			expectedResult: 1,
		},
		{
			name: "Negative #3",
			expression: "5 - (-3)",
			expectedResult: 8,
		},
		{
			name: "Negative #4",
			expression: "-2 * -3",
			expectedResult: 6,
		},

		// Плавающие числа
		{
			name: "Float #1",
			expression: "1.5 + 2.5", 
			expectedResult: 4,
		},
		{
			name: "Float #2",
			expression: "10.5 / 2.1", 
			expectedResult: 5,
		},
		{
			name: "Float #3",
			expression: "3.14 * 2", 
			expectedResult: 6.28,
		},

		//Более сложные выражения
		{
			name: "Hard #1",
			expression: "1 + 2 * (3 - 1) / 2", 
			expectedResult: 3,
		},
		{
			name: "Hard #2",
			expression: "10 / (2 + 3) * 4 - 1", 
			expectedResult: 7,
		},
		{
			name: "Hard #3",
			expression: "-(2 + 3) * 2 + 10 / 2",
			expectedResult: -5},
		{
			name: "Hard #4",
			expression: "5.2 * (3.14 - 1.14) / 2", 
			expectedResult: 5.2,
		},
		{
			name: "Hard #5",
			expression: "10/2 + 5 * 2 - (3-1)",
			expectedResult: 13,
		},
	}

	testCasesFail := []struct {
		name 				string
		expression  string
		expectedErr error
	}{
		// Синтаксические ошибки
		{
			name: "Syntax #1",
			expression: "1+1*",
			expectedErr: calculation.ErrTooManyOperats,
		},
		{
			name: "Syntax #2",
			expression: "++1", 
			expectedErr: calculation.ErrTooManyOperats,
		},
		{	name: "Syntax #3",
			expression: "(", 
			expectedErr: calculation.ErrUnclosedParen,
		},
		{
			name: "Syntax #4",
			expression: ")", 
			expectedErr: calculation.ErrUnopenedParen,
		},
		{
			name: "Syntax #5",
			expression: "1 + (2 * 3", 
			expectedErr: calculation.ErrUnclosedParen,
		},
		{
			name: "Syntax #6",
			expression: "1 + 2) * 3", 
			expectedErr: calculation.ErrUnopenedParen,
		},
		{
			name: "Syntax #7",
			expression: "1 + 2 * 3)", 
			expectedErr: calculation.ErrUnopenedParen,
		},
		{
			name: "Syntax #8",
			expression: "1.2..3", 
			expectedErr: calculation.ErrParsing,
		},
		{
			name: "Syntax #9",
			expression: "1 + + 2", 
			expectedErr: calculation.ErrTooManyOperats,
		},
		{
			name: "Syntax #10",
			expression: "1 * * 2", 
			expectedErr: calculation.ErrTooManyOperats,
		},
		{
			name: "Syntax #11",
			expression: "1 / / 2", 
			expectedErr: calculation.ErrTooManyOperats,
		},

		// Ошибки математических операций
		{
			name: "Mistakes",
			expression: "10 / 0", 
			expectedErr: calculation.ErrDivisionByZero,
		},

		// Ошибки ввода
		{
			name: "Input #1",
			expression: "abc", 
			expectedErr: calculation.ErrParsing,
		},
		{
			name: "Input #2",
			expression: "",
			expectedErr: calculation.ErrEmptyInput,
		},
		{
			name: "Input #3",
			expression: "1 + a", 
			expectedErr: calculation.ErrParsing,
		},
	}

	// Отдельное тестирование кейсов где не должна появиится ошиюка
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			// Вычисление ответа функции
			val, err := calculation.Calc(testCase.expression)
			// Предупреждаем если функция выдает ошибку
			if err != nil {
				t.Fatalf("case %s returns error", testCase.expression)
			}
			// Предупреждаем если результат вычисления неверный
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	// Отдельное тестирование кейсов где должна появиится ошиюка
	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculation.Calc(testCase.expression)
			// Предупреждаем если функция не выдает ошибку
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
			// Предупреждаем если функция выдает не ту ошибку
			if err != testCase.expectedErr {
				t.Fatalf("case %s should return error\n\t%s but got\n\t%s", testCase.expression, testCase.expectedErr, err.Error())
			}
		})
	}
}