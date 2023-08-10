package tracer

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

const traceId string = "traceId"

type contextKey struct {
}

var (
	_        slog.Handler = &Tracer{}
	traceKey              = contextKey{}
)

func MustTrace(ctx context.Context) context.Context {
	if v := ctx.Value(traceKey); v != nil {
		panic("unable to add trace to context: context already has a traceId")
	}
	uuidV4, err := uuid.NewRandom()
	if err != nil {
		panic("failed to add trace to context: " + err.Error())
	}
	return context.WithValue(ctx, traceKey, uuidV4.String())
}

type Tracer struct {
	nested slog.Handler
}

func New(nested slog.Handler) *Tracer {
	return &Tracer{nested: nested}
}

func (t *Tracer) Enabled(ctx context.Context, level slog.Level) bool {
	return t.nested.Enabled(ctx, level)
}

func (t *Tracer) Handle(ctx context.Context, record slog.Record) error {
	traceIdVal := ctx.Value(traceKey)
	if traceIdVal != nil {
		record.Add(slog.Any(traceId, traceIdVal))
	}
	return t.nested.Handle(ctx, record)
}

func (t *Tracer) WithAttrs(attrs []slog.Attr) slog.Handler {
	t.nested = t.nested.WithAttrs(attrs)
	return t.nested
}

func (t *Tracer) WithGroup(name string) slog.Handler {
	t.nested = t.nested.WithGroup(name)
	return t
}
