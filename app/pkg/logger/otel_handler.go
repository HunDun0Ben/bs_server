package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// OTelHandler 负责从 context 中提取 TraceID/SpanID 并注入到 Record 的 Attrs 中
type OTelHandler struct {
	slog.Handler
}

func NewOTelHandler(h slog.Handler) *OTelHandler {
	return &OTelHandler{Handler: h}
}

func (h *OTelHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// 注入 ID 字段，方便 Stdout JSON 能够直接显示
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()
		r.AddAttrs(
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func (h *OTelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewOTelHandler(h.Handler.WithAttrs(attrs))
}

func (h *OTelHandler) WithGroup(name string) slog.Handler {
	return NewOTelHandler(h.Handler.WithGroup(name))
}

// MultiHandler 负责将一条日志分发给多个 Handler (例如 Stdout 和 OTel Collector)
type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// 只要有一个 Handler 开启了该级别，就认为 Enabled
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(newHandlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(newHandlers...)
}
