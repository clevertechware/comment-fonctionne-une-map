// Ce fichier reproduit tel quel les deux extraits montrés dans l'article
// « Quand la map grossit », pour garantir qu'ils compilent.
package article2

import (
	"fmt"
	"maps"
	"slices"
)

func rangeOverMapHasNoGuaranteedOrder(m map[string]int) {
	for k := range m {
		fmt.Println(k) // l'ordre peut différer d'une boucle à l'autre
	}
}

func rangeOverSortedKeysIsStable(m map[string]int) {
	for _, k := range slices.Sorted(maps.Keys(m)) { // Go 1.23+
		fmt.Println(k, m[k])
	}
}
