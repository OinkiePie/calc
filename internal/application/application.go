package application

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/OinkiePie/calc/pkg/calculation"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
}

	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения в терминали
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("Введите выражение")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Не удалось прочитать выражение из консоли")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("Приложение было успешно закрыто")
			return nil
		}
		//вычисляем выражение
		result, err := calculation.Calc(text)
		
		if err != nil {
			log.Printf("Вычисление \"%s\" провалилось с ошибкой:\n\t %v", text, err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

// Серверная часть

type Request struct {
	Expression string 	`json:"expression"`
}

type Response struct {
	Status int					`json:"status"`
	Content float64			`json:"content"`
	Error string				`json:"error"`
	Timestamp int64			`json:"timestamp"`
}

type ctxReq struct{}

// Middleware для проверки валидности запроса
func RequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка типа запроса
		if r.Method != http.MethodPost {
			sendErrorResponse(w, http.StatusMethodNotAllowed, ErrOnlyPostAllowed)
			return
		}
		// Чтение и разбор запроса
		var req Request
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			// При неудачной попытки преобразования запроса отправляем оишбку с кодом 500
			sendErrorResponse(w, http.StatusInternalServerError, ErrFailedToUnmarshal)
			return
		}

		// Валидация запроса
		if req.Expression == "" {
			//  Если запрос пустой отправляем оишбку с кодом 400
			sendErrorResponse(w, http.StatusBadRequest, ErrEmptyRequest)
			return
		}

		if matched, _ := regexp.MatchString(`^[0-9/*-+() ]*$`, req.Expression); !matched {
			// Если в запросе что то кроме цифр и операторов отправляем ошибку с кодом 400

			sendErrorResponse(w, http.StatusBadRequest, ErrInvalidChars)
			return 
		}

		// Передача запроса в следующий обработчик через контекст
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxReq{}, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// sendErrorResponse - функция для отправки ответа об ошибке
func sendErrorResponse(w http.ResponseWriter, status int, err error) {
	errorResponse := Response{
		Status:    status,
		Error:     err.Error(),
		Timestamp: time.Now().Unix(),
	}

	jsonResponse, _ := json.Marshal(errorResponse)
		fmt.Fprint(w, string(jsonResponse))
}


// sendSuccessResponse - функция для отправки успешного ответа
func sendSuccessResponse(w http.ResponseWriter, answer float64) {
	succesResponse := Response{
		Status:    http.StatusOK,
		Content:   answer,
		Timestamp: time.Now().Unix(),
	}

	jsonResponse, _ := json.Marshal(succesResponse)
	fmt.Fprint(w, string(jsonResponse))
}

// CalcHandler - обработчик запроса
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	//Извлечения запроса из контекста
	request := r.Context().Value(ctxReq{}).(Request) 
	result, err := calculation.Calc(request.Expression)

	if err != nil {
		// При появлении ошибки во время вычисленя отправляем запрос с ошибкой и кодом 400
		sendErrorResponse(w, http.StatusBadRequest, err)
	} else {
		// Инчае отправляем ответ с кодом 200
		sendSuccessResponse(w, result)
	}
}

// Функция для запуска сервера
func (a *Application) RunServer() error {
	http.HandleFunc("/", RequestMiddleware(CalcHandler))
	log.Printf("Сервер был запущен на порте :%s\n", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}