package armparser

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedParserReturnsSameInstance(t *testing.T) {
	p1 := SharedParser()
	p2 := SharedParser()
	assert.Same(t, p1, p2, "SharedParser should return the same instance")
}

func TestSharedParserConcurrent(t *testing.T) {
	// Reset the shared parser to test concurrent initialization.
	parserOnce = sync.Once{}
	sharedParser = nil

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := SharedParser()
			assert.NotNil(t, p)
		}()
	}
	wg.Wait()
}
