package slogerror

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Handler is a slog.Handler that
type Handler struct {
	target slog.Handler
	attrs  []slog.Attr
}

// NewHandler creates a new *slogerror.Handler to be
// passed to slog.New(). The handler is essentially
// middleware so it requires a target handler to provide
// a final destination for log messages.
func NewHandler(target slog.Handler) *Handler {
	return &Handler{target: target}
}

// Enabled delegates the decision to the target handler
func (lh *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return lh.target.Enabled(ctx, level)
}

// Handle forward the request to the target handler
func (lh *Handler) Handle(ctx context.Context, record slog.Record) error {
	return lh.target.Handle(ctx, record)
}

// WithAttrs returns a new LogHandler with the provided attributes.
func (lh *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newTarget := lh.target.WithAttrs(attrs)
	newAttrs := lh.attrs
	if len(newAttrs) > 0 && newAttrs[len(lh.attrs)-1].Value.Kind() == slog.KindGroup {
		gAttrs := append(newAttrs[len(lh.attrs)-1].Value.Group(), attrs...)
		nGroup := slog.Attr{Key: newAttrs[len(lh.attrs)-1].Key, Value: slog.GroupValue(gAttrs...)}
		newAttrs[len(lh.attrs)-1] = nGroup
	} else {
		newAttrs = append(newAttrs, attrs...)
	}
	return &Handler{
		target: newTarget,
		attrs:  newAttrs,
	}
}

// WithGroup returns a new LogHandler with the provided group.
func (lh *Handler) WithGroup(name string) slog.Handler {
	newTarget := lh.target.WithGroup(name)
	return &Handler{
		target: newTarget,
		attrs:  append(lh.attrs, slog.Group(name)),
	}
}

func (lh *Handler) contextString(buf []byte, prefix string, a slog.Attr) []byte {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return buf
	}
	switch a.Value.Kind() {
	case slog.KindString:
		buf = fmt.Appendf(buf, "[%q = %q]", prefix+a.Key, a.Value.String())
	case slog.KindTime:
		buf = fmt.Appendf(buf, "[%q: %q]", prefix+a.Key, a.Value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, ga := range attrs {
			buf = lh.contextString(buf, prefix, ga)
		}
	default:
		buf = fmt.Appendf(buf, "[%q: %s]", prefix+a.Key, a.Value)
	}
	return buf
}
