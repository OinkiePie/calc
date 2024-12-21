package application_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OinkiePie/calc/internal/application"
)

type Response struct {
	Result float64			`json:"result"`
	Error string				`json:"error"`
	Status int					`json:"status"`
}

// Проверка сервера, получаемого статуса и ошибок.
// Проверка правильности вычисление находится в calculation_test.go
func TestServer(t *testing.T) {

	testCases := []struct {
		name           	string
		method				 	string
		request        	string
		expectedStatus 	int
		expectedError  	string
		expectedContent float64
	}{
		{
			name:           "Успешный запрос 1+1",
			method:					"POST",
			request:       `{"expression": "1 + 1"}`,
			expectedStatus: http.StatusOK,
			expectedContent: 2.0,
		},
		{
			name:           "Успешный запрос 10/2",
			method:					"POST",
			request:       `{"expression": "10 / 2"}`,
			expectedStatus: http.StatusOK,
			expectedContent: 5.0,
		},
		{
			name:           "Некорректное выражение",
			method:					"POST",
			request:       `{"expression": "invalid"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  application.ErrInvalidChars.Error(),
		},
		{
			name:           "Пустой запрос",
			method:					"POST",
			request:       `{"expression": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  application.ErrEmptyRequest.Error(),
		},
		{
			name:           "Некорректный JSON",
			method:					"POST",
			request:       `{ "expression": 5}`,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  application.ErrFailedToUnmarshal.Error(),
		},
		{
			name:           "Отсутствует поле",
			method:					"POST",
			request:       `{ "abc": "def"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  application.ErrEmptyRequest.Error(),
		},
		{
			name:           "Некорректный запрос",
			method:					"GET",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  application.ErrOnlyPostAllowed.Error(),
		},
	}

	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Создание запроса, установка header и body
			req := httptest.NewRequest(testCase.method, "/", bytes.NewBuffer([]byte(testCase.request)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder() // Создаем рекордер
			application.RequestMiddleware(http.HandlerFunc(application.CalcHandler)).ServeHTTP(w, req) // Выполняем middleware
			resp := w.Result() // Получаем ответ
			defer resp.Body.Close() // Закрываем тело овтета

			byteResponse, _ := io.ReadAll(resp.Body) // Читает тело ответа
			// Распаковывает JSON
			var response Response
			_ = json.Unmarshal(byteResponse, &response)


			if response.Status != testCase.expectedStatus { // Проверяет статус
				t.Errorf("Wrong status code for %s: expected %d, got %d", testCase.name, testCase.expectedStatus, response.Status)
			}
			if response.Error != testCase.expectedError { // Проверяет ошибку
				t.Errorf("Wrong response error for %s: expected %s, got %s", testCase.name, testCase.expectedError, response.Error)
			}
		})
	}
}

