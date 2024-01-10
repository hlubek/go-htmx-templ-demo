

## Steps

Set up a new project with Devbox and install dependencies:

```bash
git init
go mod init go-htmx-templ-todo-app
devbox init
devbox add go@1.21
devbox add tailwindcss
devbox generate direnv
tailwindcss init
```

Add templ (via Nix flake):

```bash
devbox add github:a-h/templ/v0.2.513
go get github.com/a-h/templ@v0.2.513
```

Create very basic templ app:

* `components/hello.templ`
* `main.go`

```bash
templ generate
go run main.go
```

Start a simple HTTP server that responds with a component:

* `main.go`

```
❯ curl http://localhost:3000
<div>Hello, Jane</div>%
```

Nice, now let's do some HTMX magic:

* Implement global / session counts with updates via HTMX
* `components/Layout.templ`
* `pages/Counts.templ`

Let's make it nicer!

* Use props types for components
* Extract private components
