// BenchmarkInsert1000 compare hmap (Go <= 1.23, ou GOEXPERIMENT=noswissmap sur 1.24/1.25) aux Swiss Tables (défaut
// depuis Go 1.24, seule implémentation depuis Go 1.26). Voir le README pour le relancer sous deux toolchains.
package article2

import "testing"

func BenchmarkInsert1000(b *testing.B) {
	for range b.N {
		m := make(map[int]string, 1000)
		for k := range 1000 {
			m[k] = "value"
		}
	}
}
