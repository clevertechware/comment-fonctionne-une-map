// Package article3 reproduit les extraits Mutex/RWMutex/sync.Map de l'article
// « Maps et goroutines : pourquoi ça crashe et que choisir », et fournit le
// benchmark comparatif cité dans l'article.
package article3

import "sync"

type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m == nil {
		c.m = make(map[string]int)
	}
	c.m[key]++
}
