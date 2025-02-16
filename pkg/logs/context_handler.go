package logs

import (
	"context"
	"log/slog"
)

type ctxkey string

const slogFields ctxkey = "slog_fields"

type ContextHandler struct {
	Handler slog.Handler
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ContextHandler) WithGroup(group string) slog.Handler {
	return &ContextHandler{
		Handler: h.Handler.WithGroup(group),
	}
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return h.Handler.Handle(ctx, r)
}

func AppendCtx(parent context.Context, args ...any) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	attr := argsToAttrSlice(args)

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr...)
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr...)
	return context.WithValue(parent, slogFields, v)
}

const badKey = "BAD_KEY"

func argsToAttrs(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]
	case slog.Attr:
		return x, args[1:]
	default:
		return slog.Any(badKey, x), args[1:]
	}
}

func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttrs(args)
		attrs = append(attrs, attr)
	}
	return attrs
}
