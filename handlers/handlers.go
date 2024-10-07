package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexeybs90/go_final_project/store"
)

func NextDate(res http.ResponseWriter, req *http.Request) {
	now, err := time.Parse("20060102", req.FormValue("now"))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	task := store.Task{
		Title:  "test",
		Date:   req.FormValue("date"),
		Repeat: req.FormValue("repeat"),
	}
	result, err := task.NextDate(now)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(result))
}

func SignIn(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if req.Method != http.MethodPost {
		badRequest(res, errors.New("неверный http метод"))
		return
	}
	js := JsonPass{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		okWithError(res, err)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &js); err != nil {
		okWithError(res, err)
		return
	}

	pass := js.Password
	dbpass := os.Getenv("TODO_PASSWORD")
	// log.Println(dbpass)
	if pass != dbpass {
		okWithError(res, errors.New("неверный пароль"))
		return
	}
	token, err := genToken(pass)
	if err != nil {
		okWithError(res, err)
		return
	}
	response, err := json.Marshal(JsonPass{Token: token})
	if err != nil {
		okWithError(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	// log.Println(token)
	res.Write(response)
}

func Done(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if req.Method != http.MethodPost {
		badRequest(res, http.ErrNotSupported)
	}
	id := req.URL.Query().Get("id")
	if id == "" {
		okWithError(res, errors.New("не указан id"))
		return
	}
	task_id, err := strconv.Atoi(id)
	if err != nil {
		okWithError(res, err)
		return
	}
	task, err := store.DB.Get(task_id)
	if err != nil {
		okWithError(res, err)
		return
	}
	if task.Repeat == "" {
		err = store.DB.Delete(task)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
		return
	}
	todayStr := time.Now().Format("20060102")
	today, _ := time.Parse("20060102", todayStr)
	date, err := task.NextDate(today)
	if err != nil {
		okWithError(res, err)
		return
	}
	task.Date = date
	err = store.DB.Update(task)
	if err != nil {
		okWithError(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("{}"))
}

func DoTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodPost:
		task, err := parseAndCheckTask(req)
		if err != nil {
			okWithError(res, err)
			return
		}

		id, err := store.DB.Add(task)
		if err != nil {
			okWithError(res, err)
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(JsonResult{ID: strconv.Itoa(id)})
		// res.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id)))
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		if id == "" {
			okWithError(res, errors.New("не указан id"))
			return
		}
		task_id, err := strconv.Atoi(id)
		if err != nil {
			okWithError(res, err)
			return
		}
		task, err := store.DB.Get(task_id)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(task)
	case http.MethodPut:
		task, err := parseAndCheckTask(req)
		if err != nil {
			okWithError(res, err)
			return
		}

		err = store.DB.Update(task)
		if err != nil {
			okWithError(res, err)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
	case http.MethodDelete:
		id := req.URL.Query().Get("id")
		if id == "" {
			okWithError(res, errors.New("не указан id"))
			return
		}
		task_id, err := strconv.Atoi(id)
		if err != nil {
			okWithError(res, err)
			return
		}
		task, err := store.DB.Get(task_id)
		if err != nil {
			okWithError(res, err)
			return
		}
		err = store.DB.Delete(task)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
	}
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if req.Method != http.MethodGet {
		badRequest(res, http.ErrNotSupported)
	}
	search := req.FormValue("search")
	tasks, err := store.DB.GetTasks(search)
	if err != nil {
		okWithError(res, err)
		return
	}
	json.NewEncoder(res).Encode(JsonTasks{Tasks: tasks})
}
