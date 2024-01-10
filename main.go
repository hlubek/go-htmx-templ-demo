package main

import (
	"fmt"
	"go-htmx-templ-todo-app/components"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := components.Hello("Jane")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
