package trace

import (
	"context"
	"io"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func NewOTelProvider(w io.Writer) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	return tp, tp.Shutdown, nil
}

type OTelSink struct {
	tracer   oteltrace.Tracer
	rootCtx  context.Context
	rootSpan oteltrace.Span
	mu       sync.Mutex
}

func NewOTelSink(ctx context.Context, tp *sdktrace.TracerProvider, operationName string) (*OTelSink, func()) {
	tracer := tp.Tracer("ship", oteltrace.WithInstrumentationVersion("1.0.0"))
	rootCtx, rootSpan := tracer.Start(ctx, "ship/"+operationName)
	s := &OTelSink{
		tracer:   tracer,
		rootCtx:  rootCtx,
		rootSpan: rootSpan,
	}
	return s, func() { rootSpan.End() }
}

func (s *OTelSink) Emit(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	start := now.Add(-e.Duration)

	_, span := s.tracer.Start(s.rootCtx, e.Name, oteltrace.WithTimestamp(start))
	span.SetAttributes(
		attribute.String("step", e.Name),
		attribute.Float64("duration_ms", float64(e.Duration.Microseconds())/1000.0),
	)
	if e.Err != nil {
		span.RecordError(e.Err)
		span.SetStatus(codes.Error, e.Err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End(oteltrace.WithTimestamp(now))
}
