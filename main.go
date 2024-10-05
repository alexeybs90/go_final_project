package main

import (
	"net/http"
	"os"

	"github.com/alexeybs90/go_final_project/handlers"
	"github.com/alexeybs90/go_final_project/store"
)

func main() {
	store.NewStore()
	defer store.DB.Close()

	// d, _ := time.Parse("20060102", "20240925")
	// date, err := store.NextDate(d, "20240920", "d 4")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(date)
	// return
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.HandleFunc("/api/signin", handlers.SignIn)
	http.HandleFunc("/api/task/done", handlers.Auth(handlers.Done))
	http.HandleFunc("/api/nextdate", handlers.NextDate)
	http.HandleFunc("/api/task", handlers.Auth(handlers.DoTask))
	http.HandleFunc("/api/tasks", handlers.Auth(handlers.GetTasks))
	http.Handle("/", http.FileServer(http.Dir("web")))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
