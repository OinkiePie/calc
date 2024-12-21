package calculation

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)


func Calc(expression string) (float64, error) {
	// Удаляет пробелы (Проэтому "11 11 + x" = "1111+x")
	expression = strings.ReplaceAll(expression, " ", "") // Удаление пробелов
	
	// Пустой ввод - ошибка
	if expression == "" {
		return 0, ErrEmptyInput
	}

	// Если выражение начинается с + то убираем
	// Сделано для сходжей работы "-x" в начале и "+x"
	if expression[0] == '+' && len(expression) > 1 && unicode.IsDigit(rune(expression[1])) {
		expression = strings.TrimPrefix(expression, string("+"))
	}

	expression = handleUnaryMinus(expression) // Обработка унарного минуса

	// Замена ")(число)" на ")*число"
	re := regexp.MustCompile(`\)\s*(\d+)`)
  expression = re.ReplaceAllString(expression, ")*$1")


    // Замена "(число)" на "*число"
  re = regexp.MustCompile(`(\d+)\s*\(`)
  expression = re.ReplaceAllString(expression, "$1*(")

	// Замена ")(" на ")*("
  re = regexp.MustCompile(`\)\s*\(`)
  expression = re.ReplaceAllString(expression, ")*(")

	for {
		// Проверка на правильное количество скобок и вычисление их значений
		openParen := strings.Index(expression, "(")
		if openParen == -1 {
			closeParen := findClosingParen(expression, -1)
			// Нет начальной но есть конечная скобка
			if closeParen != -1 {
				return 0, ErrUnopenedParen
			}
			
		break
		}

		closeParen := findClosingParen(expression, openParen)
		// Есть начальная но не конечная скобка
		if closeParen == -1 {
			return 0, ErrUnclosedParen
		}

		// Извлекает скобки
		innerExpr := expression[openParen+1 : closeParen]
		//Вычисляет их значение
		result, err := Calc(innerExpr)
		if err != nil {
			return 0, err
		}
		// Заменяет скобки на их значение
		expression = expression[:openParen] + fmt.Sprintf("%g", result) + expression[closeParen+1:]
	}

	// После уничтожения скобок находит значение выражения
	result, err := parseExpression(expression)
	if err != nil {
		return 0, err
	}
	// Округляем до тысячных
	result = math.Floor(result*1000)/1000
	return result, nil
}

func handleUnaryMinus(expression string) string {
	var result string
	for i, r := range expression {
		// Если минус - обрабатываем и добавляем обычную или измененную версию
		if r == '-' {
			// Минус унарный минус заменяется на арифметическое действие 
			// ( "x - y" - "x + -y")
			if i == 0 || !strings.ContainsRune("+*/", rune(expression[i-1])) {
				if (i-1 == -1 || rune(expression[i-1]) == '(') {
					//Если выражение начинается с - добавляется 0 ("-(x + y)" - 0-(x + y))
					result += "0"
					} 
				result += "+-" 
			} else {
				// Елси прсто действие- оставляем как есть
				result += "-"
			}
		} else {
			// Символ - не минус. Просто добавляем
			result += string(r)
		}
	}

	// в итоге "x -- y" заменяется на "x + y"
	result = strings.ReplaceAll(result, "+-+-", "+")
	return result
}


func findClosingParen(expr string, openIndex int) int {
	// Если попали суда значит 1 "("" уже есть
	count := 1
	// Проверяет есть ли ")" после openIndex
	for i := openIndex + 1; i < len(expr); i++ {
		if expr[i] == '(' {
			count++
		} else if expr[i] == ')' {
			count--
			// "(" и ")" равное количество
			if count == 0 {
				return i
			}
		}
	}
	// ")" не найдена
	return -1
}

func parseExpression(expression string) (float64, error) {
	parts := strings.Split(expression, "+") //Делим на слогаемые
	if len(parts) == 1 {
		//Операнд без сложения - передаем дальше
		return parseTerm(expression)
	}

	result, err := parseTerm(parts[0]) //Вычислям значение 1 слогаемого
	if err != nil {
		return 0, err
	}
	for _, part := range parts[1:] {
		//Проходим по всем слогаемым и складываем с первым
		if strings.HasPrefix(part, "-") {
			//Если слогаемое начинается с "-" вычитаем его
			//Благодаря handleUnaryMinus вычитание выглядит как сложение("x - y" стал "x + (-y)")
			operand, err := parseTerm(part[1:]) //Вычисляем слогаемое
			if err != nil {
				return 0, err
			}
			result -= operand
		} else {
			//Иначе прибавляем
			operand, err := parseTerm(part) //Вычисляем слогаемое
			if err != nil {
				return 0, err
			}
			result += operand
		}
	}
	return result, nil
}

func parseTerm(term string) (float64, error) {
	parts := strings.Split(term, "*") //Делим на множетели
	if len(parts) == 1 {
		//Операнд без умножения - передаем дальше
		return parseFactor(term)
	}

	result, err := parseFactor(parts[0]) //Вычислям значение 1 множетеля
	if err != nil {
		return 0, err
	}

	//Проходим по всем множетелям и перемножаем с первым
	for _, part := range parts[1:] {
		operand, err := parseFactor(part) //Вычисляем множитель
		if err != nil {
			return 0, err
		}
		result *= operand
	}
	return result, nil
}


func parseFactor(factor string) (float64, error) {
	parts := strings.Split(factor, "/") // Делим на делимые и делители

	if len(parts) == 1 {
		// Не получилось разбить на части т.к. ошибка или весь factor одно число
		if factor == "" {
			// Если значение пустое возвращаем ошибку
			// При вычислении с несколькими операторами подряд, лишними операторами и прочими ошибками также попадаем суда
			// ( "x//y", "x***y", "x*" и т.п. - кроме "x--y" и "--x" т.к. оно преобразуется в "x - (-y)" и "x" соответсвенно
			return 0, ErrTooManyOperats
		}
		num, err := strconv.ParseFloat(factor, 64) //Выисляем переданное число
		if err != nil {
			return num, ErrParsing
		}
		return num, nil
	}

	if parts[0] == "" {
		// Получилось разбть на части но делимое отсутствует
		// ("/x", "//x" и т.п. "-x", "+x" - намеренно сделаным исключениями )
		return 0, ErrTooManyOperats
	}
	result, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}
	for _, part := range parts[1:] {
		if part == "" {
			// Получилось разбть на части но делители отсутствуют
			// ("x/", "x//"  и т.п. )
			return 0, ErrTooManyOperats
		}
		operand, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return 0, err
		}
		if operand == 0 {
			// Делитель 0 - ошибка
			return 0, ErrDivisionByZero
		}
		result /= operand
	}
	return result, nil
}