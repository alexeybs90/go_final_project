package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

func DoTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodPost:
		task, err := parseAndCheckTask(req)
		if err != nil {
			badRequest(res, err)
			return
		}

		id, err := store.DB.Add(task)
		if err != nil {
			badRequest(res, err)
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(JsonResult{ID: strconv.Itoa(id)})
		// res.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id)))
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		if id == "" {
			badRequest(res, errors.New("не указан id"))
			return
		}
		task_id, err := strconv.Atoi(id)
		if err != nil {
			badRequest(res, err)
			return
		}
		task, err := store.DB.Get(task_id)
		if err != nil {
			badRequest(res, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(task)
	case http.MethodPut:
		task, err := parseAndCheckTask(req)
		if err != nil {
			badRequest(res, err)
			return
		}

		err = store.DB.Update(task)
		if err != nil {
			badRequest(res, err)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
	}
}

func parseAndCheckTask(req *http.Request) (store.Task, error) {
	//task := new(store.Task)
	task := store.Task{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, err
	}

	if req.Method == http.MethodPut {
		//проверим, есть ли такая таска в бд
		task_id, err := strconv.Atoi(task.ID)
		if err != nil {
			return task, err
		}
		_, err = store.DB.Get(task_id)
		if err != nil {
			return task, err
		}
	}

	// log.Println(task)
	err = task.CheckSave()
	if err != nil {
		return task, err
	}

	todayStr := time.Now().Format("20060102")
	today, _ := time.Parse("20060102", todayStr) //отбрасываем время, чтобы было 00:00:00
	if task.Date == "" {
		task.Date = todayStr
	} else {
		newD, _ := time.Parse("20060102", task.Date)
		// log.Println(newD, today)
		if newD.Before(today) {
			if task.Repeat == "" {
				task.Date = today.Format("20060102")
			} else {
				task.Date, _ = task.NextDate(today)
			}
		}
	}
	return task, err
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if req.Method != http.MethodGet {
		badRequest(res, http.ErrNotSupported)
	}
	search := req.FormValue("search")
	tasks, err := store.DB.GetTasks(search)
	if err != nil {
		badRequest(res, err)
		return
	}
	json.NewEncoder(res).Encode(JsonTasks{Tasks: tasks})
}

func badRequest(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(res).Encode(JsonResult{Error: err.Error()})
	// res.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
}

type JsonResult struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type JsonTasks struct {
	Tasks []store.Task `json:"tasks"`
}
