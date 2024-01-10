package handler

import (
	"net/http"

	"github.com/gorilla/sessions"

	"go-htmx-templ-todo-app/pages"
	"go-htmx-templ-todo-app/service"
)

var store = sessions.NewCookieStore([]byte("only-for-development"))

type Handler struct {
	mux *http.ServeMux

	counter service.Counter
}

func New(counter service.Counter) *Handler {
	h := &Handler{
		mux:     http.NewServeMux(),
		counter: counter,
	}

	h.mux.HandleFunc("/", h.root)

	return h
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.create(w, r)
		return
	}
	h.get(w, r)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")

	sessionCount := 0
	if count, ok := session.Values["counter"].(int); ok {
		sessionCount = count
	}

	globalCount := h.counter.Get()

	pages.CountsPage(globalCount, sessionCount).Render(r.Context(), w)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var globalCount, sessionCount int

	// Get session count
	session, _ := store.Get(r, "user")
	sessionCount = 0
	if count, ok := session.Values["counter"].(int); ok {
		sessionCount = count
	}

	// Decide the action to take based on the button that was pressed.
	if r.Form.Has("session") {
		sessionCount++
		session.Values["counter"] = sessionCount
		session.Save(r, w)
	}
	if r.Form.Has("global") {
		globalCount = h.counter.Increment()
	} else {
		globalCount = h.counter.Get()
	}

	pages.CountsForm(globalCount, sessionCount).Render(r.Context(), w)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
