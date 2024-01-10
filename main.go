package main

import (
	"fmt"
	"net/http"

	"go-htmx-templ-todo-app/handler"
	"go-htmx-templ-todo-app/service"
)

func main() {
	counter := service.NewInMemoryCounter()

	srv := http.NewServeMux()
	srv.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))
	srv.Handle("/", handler.New(counter))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", srv)
}
