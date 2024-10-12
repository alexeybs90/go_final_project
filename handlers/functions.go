package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexeybs90/go_final_project/store"
	"github.com/golang-jwt/jwt/v5"
)

type JsonPass struct {
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type JsonResult struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type JsonTasks struct {
	Tasks []store.Task `json:"tasks"`
}

func badRequest(res http.ResponseWriter, err error) {
	serverError(res, http.StatusBadRequest, err)
}
func okWithError(res http.ResponseWriter, err error) {
	serverError(res, http.StatusOK, err)
}

func serverError(res http.ResponseWriter, status int, err error) {
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(JsonResult{Error: err.Error()})
}

func genToken(pass string) (string, error) {
	secret := []byte(pass + "_todo_secret")
	// создаём jwt и указываем алгоритм хеширования
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	// получаем подписанный токен
	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s TaskService) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := s.todoPassword
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			token, err := genToken(pass)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if jwt != token {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

// parseAndCheckTask парсит таску из json и проверяет на ошибки
func (s TaskService) parseAndCheckTask(req *http.Request) (store.Task, error) {
	task := store.Task{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, err
	}

	err = task.Validate()
	if err != nil {
		return task, err
	}

	todayStr := time.Now().Format(store.DateFormat)
	today, _ := time.Parse(store.DateFormat, todayStr) //отбрасываем время, чтобы было 00:00:00
	if task.Date == "" {
		task.Date = todayStr
	} else {
		newD, _ := time.Parse(store.DateFormat, task.Date)
		if newD.Before(today) {
			if task.Repeat == "" {
				task.Date = today.Format(store.DateFormat)
			} else {
				task.Date, _ = task.NextDate(today)
			}
		}
	}
	return task, err
}
