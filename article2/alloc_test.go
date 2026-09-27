// Package article2 mesure le coût d'une insertion de 1000 clés, avec et
// sans pré-dimensionnement (make(map[K]V, n)), pour le deuxième article de
// la série. Le nombre de buckets ou de tables allouées n'est pas exposé par
// l'API : testing.AllocsPerRun compte les allocations mémoire, qui est le
// signal observable proprement depuis l'extérieur du runtime.
package article2

import "testing"

func insertNKeys(m map[int]int, n int) {
	for i := range n {
		m[i] = i
	}
}

func TestPreSizedMapAllocatesLessThanGrowingMap(t *testing.T) {
	const n = 1000

	growing := testing.AllocsPerRun(50, func() {
		m := make(map[int]int)
		insertNKeys(m, n)
	})

	preSized := testing.AllocsPerRun(50, func() {
		m := make(map[int]int, n)
		insertNKeys(m, n)
	})

	if preSized >= growing {
		t.Errorf("pré-dimensionner devrait réduire les allocations: pré-dimensionné=%.1f, croissance=%.1f", preSized, growing)
	}
}

func BenchmarkInsertGrowing(b *testing.B) {
	for range b.N {
		m := make(map[int]int)
		insertNKeys(m, 1000)
	}
}

func BenchmarkInsertPreSized(b *testing.B) {
	for range b.N {
		m := make(map[int]int, 1000)
		insertNKeys(m, 1000)
	}
}
