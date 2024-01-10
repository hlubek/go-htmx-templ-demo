package main

import (
	"context"
	"os"

	"go-htmx-templ-todo-app/components"
)

func main() {
	component := components.Hello("John")
	component.Render(context.Background(), os.Stdout)
}
