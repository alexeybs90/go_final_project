package handlers

import (
	"bytes"
	"encoding/json"
	"log"
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

func PostTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodPost:
		//task := new(store.Task)
		task := store.Task{}
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			badRequest(res, err)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			badRequest(res, err)
			return
		}

		log.Println(task)

		err = task.CheckSave()
		if err != nil {
			badRequest(res, err)
			return
		}

		today := time.Now()
		if task.Date == "" {
			task.Date = today.Format("20060102")
		} else {
			newD, _ := time.Parse("20060102", task.Date)
			if newD.Before(today) {
				if task.Repeat == "" {
					task.Date = today.Format("20060102")
				} else {
					task.Date, _ = task.NextDate(today)
				}
			}
		}

		id, err := store.DB.Add(task)
		if err != nil {
			badRequest(res, err)
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(JsonResult{ID: strconv.Itoa(id)})
		// res.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id)))
	}
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
