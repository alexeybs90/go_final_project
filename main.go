package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexeybs90/go_final_project/handlers"
	"github.com/alexeybs90/go_final_project/store"
)

func main() {
	s := store.NewStore()
	defer s.Close()

	todoPassword := os.Getenv("TODO_PASSWORD")
	service := handlers.NewTaskService(s, todoPassword)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.HandleFunc("/api/signin", service.SignIn)
	http.HandleFunc("/api/task/done", service.Auth(service.Done))
	http.HandleFunc("/api/nextdate", service.NextDate)
	http.HandleFunc("/api/task", service.Auth(service.DoTask))
	http.HandleFunc("/api/tasks", service.Auth(service.GetTasks))
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Println("Стартуем...")

	go func() {
		time.Sleep(time.Second * 3)
		_, err := http.Get("http://localhost:" + port)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Приложение запущено на порту " + port)
	}()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
