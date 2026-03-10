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
	var wg sync.WaitGroup
	parsers := make([]*any, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p := SharedParser()
			assert.NotNil(t, p)
			v := any(p)
			parsers[idx] = &v
		}(i)
	}
	wg.Wait()
}
