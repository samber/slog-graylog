package main

import (
	"fmt"
	"log"
	"time"

	"log/slog"

	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
)

func main() {
	// docker-compose up -d
	// or
	// ncat -l 12201 -u
	gelfWriter, err := gelf.NewWriter("localhost:12201")
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}

	logger := slog.New(sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}
