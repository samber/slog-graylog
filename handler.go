package sloggraylog

import (
	"context"
	"encoding/json"

	"log/slog"

	"github.com/Graylog2/go-gelf/gelf"
	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// connection to graylog
	Writer *gelf.Writer

	// optional: customize json payload builder
	Converter Converter
	// optional: custom marshaler
	Marshaler func(v any) ([]byte, error)

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (o Option) NewGraylogHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Writer == nil {
		panic("missing graylog connections")
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.Marshaler == nil {
		o.Marshaler = json.Marshal
	}

	return &GraylogHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*GraylogHandler)(nil)

type GraylogHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *GraylogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *GraylogHandler) Handle(ctx context.Context, record slog.Record) error {
	message := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	bytes, err := h.option.Marshaler(message)
	if err != nil {
		return err
	}

	go func() {
		_, _ = h.option.Writer.Write(append(bytes, byte('\n')))
	}()

	return nil
}

func (h *GraylogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GraylogHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
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
