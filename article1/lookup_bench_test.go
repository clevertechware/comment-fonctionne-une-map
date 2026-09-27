// Package article1 mesure le coût d'un lookup réussi et d'un lookup manqué
// sur une map[string]int, pour le premier article de la série. Lancer avec
// -benchmem pour voir le nombre d'allocations par opération.
package article1

import (
	"fmt"
	"testing"
)

func buildStringIntMap(n int) map[string]int {
	m := make(map[string]int, n)
	for i := range n {
		m[fmt.Sprintf("key-%04d", i)] = i
	}
	return m
}

func BenchmarkLookupHit(b *testing.B) {
	m := buildStringIntMap(1000)
	key := "key-0500"

	for range b.N {
		_ = m[key]
	}
}

func BenchmarkLookupMiss(b *testing.B) {
	m := buildStringIntMap(1000)
	key := "key-9999"

	for range b.N {
		_ = m[key]
	}
}
