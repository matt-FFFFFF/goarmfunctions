package armparser

import (
	"sync"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/stretchr/testify/assert"
)

func TestSharedParserReturnsSameInstance(t *testing.T) {
	p1 := SharedParser()
	p2 := SharedParser()
	assert.Same(t, p1, p2, "SharedParser should return the same instance")
}

func TestSharedParserConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	parsers := make([]*participle.Parser[ArmValue], 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p := SharedParser()
			assert.NotNil(t, p)
			parsers[idx] = p
		}(i)
	}
	wg.Wait()

	// Verify that all goroutines received the same shared parser instance.
	first := parsers[0]
	assert.NotNil(t, first)
	for i := 1; i < len(parsers); i++ {
		assert.NotNil(t, parsers[i])
		assert.Same(t, first, parsers[i], "all goroutines should see the same SharedParser instance")
	}
}
