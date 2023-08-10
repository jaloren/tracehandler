package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jaloren/tracehandler/tracer"
	"log/slog"
	"os"
)

func init() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	traceHandler := tracer.New(jsonHandler)
	logger := slog.New(traceHandler)
	slog.SetDefault(logger)
}

func main() {
	ctx := tracer.MustTrace(context.Background())
	if err := processPipeline(ctx); err != nil {
		slog.Error("failed to process the pipeline")
	}
}

func processPipeline(ctx context.Context) error {
	return transformOne(ctx)
}

func transformOne(ctx context.Context) error {
	err := transformTwo(ctx)
	slog.ErrorContext(ctx, "transform one failed")
	return fmt.Errorf("transform one failed: %w", err)

}

func transformTwo(ctx context.Context) error {
	err := transformThree(ctx)
	slog.ErrorContext(ctx, "second transform failed")
	return fmt.Errorf("second transform failed: %w", err)
}

func transformThree(ctx context.Context) error {
	multiErr := []error{errors.New("first error"), errors.New("second error")}
	err := errors.Join(multiErr...)
	slog.ErrorContext(ctx, "third transform failed", "err", err)
	return fmt.Errorf("third transform failed: %w", err)
}
