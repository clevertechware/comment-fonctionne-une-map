// BenchmarkInsert1000 sert à comparer l'implémentation hmap (Go <= 1.23,
// GOEXPERIMENT=noswissmap sur 1.24/1.25) aux Swiss Tables (par défaut depuis
// Go 1.24, seule implémentation depuis Go 1.26). Voir le README pour lancer
// ce même bench avec deux toolchains différentes : ce fichier ne fabrique
// aucun chiffre de comparaison, il ne fait que fournir le bench à relancer.
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
