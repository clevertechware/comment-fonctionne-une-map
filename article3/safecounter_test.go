package article3

import (
	"sync"
	"testing"
)

func TestSafeCounterZeroValueIncDoesNotPanic(t *testing.T) {
	var c SafeCounter

	c.Inc("shared-key")

	if got, want := c.m["shared-key"], 1; got != want {
		t.Errorf("compteur = %d, voulu %d", got, want)
	}
}

func TestSafeCounterIncIsRaceFreeUnderConcurrentWrites(t *testing.T) {
	var c SafeCounter

	var wg sync.WaitGroup
	const goroutines, incrementsPerGoroutine = 50, 100
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range incrementsPerGoroutine {
				c.Inc("shared-key")
			}
		}()
	}
	wg.Wait()

	want := goroutines * incrementsPerGoroutine
	if got := c.m["shared-key"]; got != want {
		t.Errorf("compteur = %d, voulu %d", got, want)
	}
}
