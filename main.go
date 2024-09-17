package main

import (
	"net/http"
	"os"

	"github.com/alexeybs90/go_final_project/store"
)

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("hi"))
}

func main() {
	st := store.NewStore()
	defer st.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.HandleFunc("/tasks", ServeHTTP)
	http.Handle("/", http.FileServer(http.Dir("web")))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
