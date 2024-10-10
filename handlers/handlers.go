package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexeybs90/go_final_project/store"
)

type TaskService struct {
	store        store.Store
	TodoPassword string
}

func NewTaskService(store store.Store) TaskService {
	return TaskService{store: store}
}

func (s TaskService) NextDate(res http.ResponseWriter, req *http.Request) {
	now, err := time.Parse(store.DateFormat, req.FormValue("now"))
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
	_, err = res.Write([]byte(result))
	if err != nil {
		log.Println(err.Error())
	}
}

func (s TaskService) SignIn(res http.ResponseWriter, req *http.Request) {
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
	if pass != s.TodoPassword {
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
	_, err = res.Write(response)
	if err != nil {
		log.Println(err.Error())
	}
}

func (s TaskService) Done(res http.ResponseWriter, req *http.Request) {
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
	task, err := s.store.Get(task_id)
	if err != nil {
		okWithError(res, err)
		return
	}
	if task.Repeat == "" {
		err = s.store.Delete(task)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("{}"))
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	todayStr := time.Now().Format(store.DateFormat)
	today, _ := time.Parse(store.DateFormat, todayStr)
	date, err := task.NextDate(today)
	if err != nil {
		okWithError(res, err)
		return
	}
	task.Date = date
	err = s.store.Update(task)
	if err != nil {
		okWithError(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte("{}"))
	if err != nil {
		log.Println(err.Error())
	}
}

func (s TaskService) DoTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodPost:
		task, err := s.parseAndCheckTask(req)
		if err != nil {
			okWithError(res, err)
			return
		}

		id, err := s.store.Add(task)
		if err != nil {
			okWithError(res, err)
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(JsonResult{ID: strconv.Itoa(id)})
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
		task, err := s.store.Get(task_id)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(task)
	case http.MethodPut:
		task, err := s.parseAndCheckTask(req)
		if err != nil {
			okWithError(res, err)
			return
		}

		err = s.store.Update(task)
		if err != nil {
			okWithError(res, err)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("{}"))
		if err != nil {
			log.Println(err.Error())
		}
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
		task, err := s.store.Get(task_id)
		if err != nil {
			okWithError(res, err)
			return
		}
		err = s.store.Delete(task)
		if err != nil {
			okWithError(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("{}"))
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (s TaskService) GetTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if req.Method != http.MethodGet {
		badRequest(res, http.ErrNotSupported)
	}
	search := req.FormValue("search")
	tasks, err := s.store.GetTasks(search)
	if err != nil {
		okWithError(res, err)
		return
	}
	json.NewEncoder(res).Encode(JsonTasks{Tasks: tasks})
}
