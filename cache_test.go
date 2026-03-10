package goarmfunctions

import (
	"context"
	"sync"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/armparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateCacheHit(t *testing.T) {
	// Clear cache before test.
	parseCache = sync.Map{}

	expr := "[if(equals('a', 'a'), 'yes', 'no')]"
	ctx := context.Background()

	// First call should parse and cache.
	result1, err := Evaluate(ctx, expr, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "yes", result1)

	// Verify the expression is cached.
	_, ok := parseCache.Load(expr)
	assert.True(t, ok, "expression should be cached after first call")

	// Second call should use cache and produce the same result.
	result2, err := Evaluate(ctx, expr, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "yes", result2)
}

func TestEvaluateCacheDifferentContexts(t *testing.T) {
	// Clear cache before test.
	parseCache = sync.Map{}

	expr := "[parameters('key')]"
	ctx := context.Background()

	// Evaluate with first context.
	evalCtx1 := armparser.FromMap(map[string]any{"key": "value1"})
	result1, err := Evaluate(ctx, expr, evalCtx1, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "value1", result1)

	// Evaluate with different context using the same (cached) expression.
	evalCtx2 := armparser.FromMap(map[string]any{"key": "value2"})
	result2, err := Evaluate(ctx, expr, evalCtx2, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "value2", result2)
}

func TestEvaluateCacheConcurrent(t *testing.T) {
	// Clear cache before test.
	parseCache = sync.Map{}

	expr := "[if(equals('a', 'a'), 'yes', 'no')]"
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := Evaluate(ctx, expr, nil, nil, nil)
			assert.NoError(t, err)
			assert.Equal(t, "yes", result)
		}()
	}
	wg.Wait()
}

func TestEvaluateParseError(t *testing.T) {
	// Clear cache before test.
	parseCache = sync.Map{}

	expr := "[invalid expression!!!"
	ctx := context.Background()

	_, err := Evaluate(ctx, expr, nil, nil, nil)
	require.Error(t, err)

	// Verify that parse errors are not cached.
	_, ok := parseCache.Load(expr)
	assert.False(t, ok, "parse errors should not be cached")
}

func BenchmarkEvaluateNoCache(b *testing.B) {
	expr := "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]"
	ctx := context.Background()
	registry := armparser.DefaultRegistry()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Clear cache to simulate no caching.
		parseCache = sync.Map{}
		_, err := Evaluate(ctx, expr, nil, registry, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvaluateWithCache(b *testing.B) {
	expr := "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]"
	ctx := context.Background()
	registry := armparser.DefaultRegistry()

	// Prime the cache.
	parseCache = sync.Map{}
	_, err := Evaluate(ctx, expr, nil, registry, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Evaluate(ctx, expr, nil, registry, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
