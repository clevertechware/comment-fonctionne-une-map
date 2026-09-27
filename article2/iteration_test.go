package article2

import (
	"fmt"
	"maps"
	"slices"
	"testing"
)

func collectIterationOrder(m map[string]int) []string {
	var order []string
	for k := range m {
		order = append(order, k)
	}
	return order
}

func TestTwoRangesOverSameMapCanProduceDifferentOrders(t *testing.T) {
	m := make(map[string]int, 64)
	for i := range 64 {
		m[fmt.Sprintf("k%02d", i)] = i
	}

	sameEveryTime := true
	first := collectIterationOrder(m)
	for range 20 {
		if !slices.Equal(collectIterationOrder(m), first) {
			sameEveryTime = false
			break
		}
	}

	if sameEveryTime {
		t.Error("ordre identique sur 20 tirages avec 64 clés : statistiquement improbable, l'itération n'est peut-être plus randomisée")
	}
}

func TestSortedKeysGivesStableOrder(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}

	first := slices.Sorted(maps.Keys(m))
	second := slices.Sorted(maps.Keys(m))

	if !slices.Equal(first, second) {
		t.Fatalf("ordre trié instable: %v puis %v", first, second)
	}
	if !slices.Equal(first, []string{"a", "b", "c"}) {
		t.Fatalf("ordre trié incorrect: %v", first)
	}
}
