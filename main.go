package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"

	"go-htmx-templ-todo-app/handler"
	"go-htmx-templ-todo-app/service"
)

func main() {
	setDefaultLogger()

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Value: ":3000",
		},
		// Bind to env vars that are set by refresh if live reload is enabled
		&cli.StringFlag{
			Name:    "live-reload-sse-url",
			EnvVars: []string{"REFRESH_LIVE_RELOAD_SSE_URL"},
		},
		&cli.StringFlag{
			Name:    "live-reload-sse-event",
			EnvVars: []string{"REFRESH_LIVE_RELOAD_SSE_EVENT"},
		},
	}
	app.Action = server
	err := app.Run(os.Args)
	if err != nil {
		slog.Error("Command failed", "err", err)
		os.Exit(1)
	}
}

func setDefaultLogger() {
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
			NoColor:    !isatty.IsTerminal(w.Fd()),
		}),
	))
}

func server(c *cli.Context) error {
	counter := service.NewInMemoryCounter()

	config := handler.Config{
		LiveReloadSSEurl:   c.String("live-reload-sse-url"),
		LiveReloadSSEevent: c.String("live-reload-sse-event"),
	}
	if c.String("live-reload-sse-url") != "" {
		slog.Debug("Live reload enabled", "liveReloadSSEurl", c.String("live-reload-sse-url"))
	}

	srv := http.NewServeMux()
	srv.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))
	srv.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	srv.Handle("/", handler.New(config, counter))

	slog.Info("Listening", "address", c.String("address"))
	return http.ListenAndServe(c.String("address"), srv)
}
