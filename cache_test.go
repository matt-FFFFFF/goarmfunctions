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
	ResetParseCache()

	expr := "[if(equals('a', 'a'), 'yes', 'no')]"
	ctx := context.Background()

	// First call should parse and cache.
	result1, err := Evaluate(ctx, expr, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "yes", result1)

	// Verify the expression is cached.
	_, ok := parseCache.load(expr)
	assert.True(t, ok, "expression should be cached after first call")

	// Second call should use cache and produce the same result.
	result2, err := Evaluate(ctx, expr, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "yes", result2)
}

func TestEvaluateCacheDifferentContexts(t *testing.T) {
	ResetParseCache()

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
	ResetParseCache()

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
	ResetParseCache()

	expr := "[invalid expression!!!"
	ctx := context.Background()

	_, err := Evaluate(ctx, expr, nil, nil, nil)
	require.Error(t, err)

	// Verify that parse errors are not cached.
	_, ok := parseCache.load(expr)
	assert.False(t, ok, "parse errors should not be cached")
}

func TestResetParseCache(t *testing.T) {
	ResetParseCache()

	expr := "[if(equals('a', 'a'), 'yes', 'no')]"
	ctx := context.Background()

	_, err := Evaluate(ctx, expr, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, ParseCacheLen())

	ResetParseCache()
	assert.Equal(t, 0, ParseCacheLen())
}

func TestSetParseCacheSize(t *testing.T) {
	ResetParseCache()
	SetParseCacheSize(2)
	defer SetParseCacheSize(DefaultParseCacheSize) // restore default after test

	ctx := context.Background()
	exprs := []string{
		"[if(equals('a', 'a'), 'yes', 'no')]",
		"[if(equals('b', 'b'), 'yes', 'no')]",
		"[if(equals('c', 'c'), 'yes', 'no')]",
	}

	// Fill cache to capacity.
	for _, expr := range exprs[:2] {
		_, err := Evaluate(ctx, expr, nil, nil, nil)
		require.NoError(t, err)
	}
	assert.Equal(t, 2, ParseCacheLen())

	// Adding a third entry should evict the LRU (first) entry.
	_, err := Evaluate(ctx, exprs[2], nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, ParseCacheLen())

	// First expression should have been evicted.
	_, ok := parseCache.load(exprs[0])
	assert.False(t, ok, "first expression should have been evicted")

	// Second and third expressions should still be cached.
	_, ok = parseCache.load(exprs[1])
	assert.True(t, ok, "second expression should still be cached")
	_, ok = parseCache.load(exprs[2])
	assert.True(t, ok, "third expression should still be cached")
}

func TestSetParseCacheSizeShrink(t *testing.T) {
	ResetParseCache()
	SetParseCacheSize(DefaultParseCacheSize)
	defer SetParseCacheSize(DefaultParseCacheSize)

	ctx := context.Background()
	exprs := []string{
		"[if(equals('a', 'a'), 'yes', 'no')]",
		"[if(equals('b', 'b'), 'yes', 'no')]",
		"[if(equals('c', 'c'), 'yes', 'no')]",
	}

	for _, expr := range exprs {
		_, err := Evaluate(ctx, expr, nil, nil, nil)
		require.NoError(t, err)
	}
	assert.Equal(t, 3, ParseCacheLen())

	// Shrinking below current size should evict LRU entries.
	SetParseCacheSize(1)
	assert.Equal(t, 1, ParseCacheLen())

	// Only the most recently used entry should remain.
	_, ok := parseCache.load(exprs[2])
	assert.True(t, ok, "most recently used expression should still be cached")
}

func TestLRUEvictionOrder(t *testing.T) {
	ResetParseCache()
	SetParseCacheSize(2)
	defer SetParseCacheSize(DefaultParseCacheSize)

	ctx := context.Background()
	expr1 := "[if(equals('a', 'a'), 'yes', 'no')]"
	expr2 := "[if(equals('b', 'b'), 'yes', 'no')]"
	expr3 := "[if(equals('c', 'c'), 'yes', 'no')]"

	_, err := Evaluate(ctx, expr1, nil, nil, nil)
	require.NoError(t, err)
	_, err = Evaluate(ctx, expr2, nil, nil, nil)
	require.NoError(t, err)

	// Access expr1 again to make it most recently used.
	_, err = Evaluate(ctx, expr1, nil, nil, nil)
	require.NoError(t, err)

	// Adding expr3 should evict expr2 (least recently used), not expr1.
	_, err = Evaluate(ctx, expr3, nil, nil, nil)
	require.NoError(t, err)

	_, ok := parseCache.load(expr1)
	assert.True(t, ok, "expr1 should still be cached (recently accessed)")
	_, ok = parseCache.load(expr2)
	assert.False(t, ok, "expr2 should have been evicted (LRU)")
	_, ok = parseCache.load(expr3)
	assert.True(t, ok, "expr3 should be cached (just added)")
}

func BenchmarkEvaluateNoCache(b *testing.B) {
	expr := "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]"
	ctx := context.Background()
	registry := armparser.DefaultRegistry()

	b.ReportAllocs()
	for b.Loop() {
		// Clear cache to simulate no caching.
		ResetParseCache()
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
	ResetParseCache()
	_, err := Evaluate(ctx, expr, nil, registry, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, err := Evaluate(ctx, expr, nil, registry, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
