package sloggraylog

import (
	"context"
	"encoding/json"

	"github.com/Graylog2/go-gelf/gelf"
	"golang.org/x/exp/slog"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// connection to graylog
	Writer *gelf.Writer

	// optional: customize json payload builder
	Converter Converter
}

func (o Option) NewGraylogHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Writer == nil {
		panic("missing graylog connections")
	}

	return &GraylogHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

type GraylogHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *GraylogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *GraylogHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	message := converter(h.attrs, &record)

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = h.option.Writer.Write(append(bytes, byte('\n')))

	return err
}

func (h *GraylogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GraylogHandler{
		option: h.option,
		attrs:  appendAttrsToGroup(h.groups, h.attrs, attrs),
		groups: h.groups,
	}
}

func (h *GraylogHandler) WithGroup(name string) slog.Handler {
	return &GraylogHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
