package handler

import (
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"github.com/r3labs/sse/v2"

	"go-htmx-templ-todo-app/components"
	"go-htmx-templ-todo-app/pages"
	"go-htmx-templ-todo-app/service"
)

var store = sessions.NewCookieStore([]byte("only-for-development"))

type Config struct {
	LiveReloadSSEurl   string
	LiveReloadSSEevent string
}

type Handler struct {
	mux *http.ServeMux

	counter service.Counter

	s *sse.Server

	config Config
}

func New(config Config, counter service.Counter) *Handler {
	s := sse.New()
	s.AutoReplay = false
	s.CreateStream("global")

	h := &Handler{
		mux:     http.NewServeMux(),
		counter: counter,
		s:       s,

		config: config,
	}

	h.mux.HandleFunc("/sse", h.events)
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
	session := h.getOrCreateSession(w, r)

	globalCount := h.counter.Get()
	sessionCount := session.Values["counter"].(int)

	counts := pages.Counts{Global: globalCount, Session: sessionCount}
	pages.CountsPage(pages.CountsPageProps{Counts: counts, LayoutProps: h.layoutProps()}).Render(r.Context(), w)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var globalCount, sessionCount int

	// Get session count
	session := h.getOrCreateSession(w, r)
	sessionCount = session.Values["counter"].(int)

	// Decide the action to take based on the button that was pressed.
	if r.Form.Has("session") {
		sessionCount++
		session.Values["counter"] = sessionCount
		session.Save(r, w)

		h.s.Publish(session.Values["clientID"].(string), &sse.Event{
			Event: []byte("SessionCountChanged"),
			Data:  []byte(`<span>` + strconv.Itoa(sessionCount) + `</span>`),
		})
	}
	if r.Form.Has("global") {
		globalCount = h.counter.Increment()

		h.s.Publish("global", &sse.Event{
			Event: []byte("GlobalCountChanged"),
			Data:  []byte(`<span>` + strconv.Itoa(globalCount) + `</span>`),
		})
	} else {
		globalCount = h.counter.Get()
	}

	counts := pages.Counts{Global: globalCount, Session: sessionCount}
	pages.CountsForm(pages.CountsFormProps{Counts: counts}).Render(r.Context(), w)
}

func (h *Handler) getOrCreateSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "user")
	clientID := ""
	if session.Values["clientID"] == nil || session.Values["counter"] == nil {
		clientID = uuid.Must(uuid.NewV4()).String()
		session.Values["clientID"] = clientID
		session.Values["counter"] = 0
		session.Save(r, w)
	}

	return session
}

func (h *Handler) events(w http.ResponseWriter, r *http.Request) {
	stream := r.FormValue("stream")
	if stream != "global" {
		session := h.getOrCreateSession(w, r)
		clientID := session.Values["clientID"].(string)

		// Create stream for session
		if !h.s.StreamExists(clientID) {
			h.s.CreateStream(clientID)
		}

		// Serve session events stream by passing the session id as "stream" query param to the sse handler
		values := r.URL.Query()

		values.Set("stream", clientID)
		r.URL.RawQuery = values.Encode()

		h.s.ServeHTTP(w, r)
		return
	}

	h.s.ServeHTTP(w, r)
}

func (h *Handler) layoutProps() components.LayoutProps {
	return components.LayoutProps{
		LiveReloadSSEurl:   h.config.LiveReloadSSEurl,
		LiveReloadSSEevent: h.config.LiveReloadSSEevent,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
