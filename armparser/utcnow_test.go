package armparser

import (
	"context"
	"strings"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestUtcNow(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	parser := New()
	f, err := parser.ParseString("test", "[utcNow()]")
	if err != nil {
		t.Fatal(err)
	}
	result, err := f.Evaluate(ctx, FromMap(map[string]any{}), DefaultRegistry())
	if err != nil {
		t.Fatal(err)
	}
	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	// Check format: should end with Z and contain T
	if !strings.Contains(s, "T") || !strings.HasSuffix(s, "Z") {
		t.Errorf("unexpected format: %s", s)
	}
}
