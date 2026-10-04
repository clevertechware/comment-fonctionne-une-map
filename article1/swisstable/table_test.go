package swisstable

import (
	"fmt"
	"hash/maphash"
	"slices"
	"testing"
)

const maxKeysPerTable = maxTableCapacity * maxAvgGroupLoad / groupSize

func newTestTable() *Table {
	return NewTable(maphash.MakeSeed())
}

func testKey(i int) string {
	return fmt.Sprintf("cle-%d", i)
}

func TestTableDoublesItsCapacityWhenLoadExceedsSevenEighths(t *testing.T) {
	table := newTestTable()

	for i := range 7 {
		table.Insert(testKey(i))
	}
	if got := table.Capacity(); got != 8 {
		t.Fatalf("capacité après 7 clés = %d, voulu 8 (7/8 de charge n'est pas dépassé)", got)
	}

	table.Insert(testKey(7))
	if got := table.Capacity(); got != 16 {
		t.Errorf("capacité après 8 clés = %d, voulu 16", got)
	}
}

func TestTableKeepsEveryKeyAfterGrowing(t *testing.T) {
	const n = maxKeysPerTable
	table := newTestTable()

	for i := range n {
		table.Insert(testKey(i))
	}

	for i := range n {
		if !table.Contains(testKey(i)) {
			t.Fatalf("clé %q perdue pendant la croissance", testKey(i))
		}
	}
	if got := table.Len(); got != n {
		t.Errorf("Len = %d, voulu %d", got, n)
	}
}

func TestTableNeverExceedsSevenEighthsLoadAfterAnInsert(t *testing.T) {
	table := newTestTable()

	for i := range maxKeysPerTable {
		table.Insert(testKey(i))
		if table.Len()*groupSize > table.Capacity()*maxAvgGroupLoad {
			t.Fatalf("après %d clés, %d/%d slots occupés dépasse 7/8", i+1, table.Len(), table.Capacity())
		}
	}
}

func TestInsertingAnExistingKeyChangesNeitherLenNorCapacity(t *testing.T) {
	table := newTestTable()
	for i := range 7 {
		table.Insert(testKey(i))
	}

	table.Insert(testKey(0))

	if table.Len() != 7 || table.Capacity() != 8 {
		t.Errorf("Len=%d Capacity=%d après réinsertion, voulu 7 et 8", table.Len(), table.Capacity())
	}
}

func TestContainsRejectsAbsentKeyAfterGrowing(t *testing.T) {
	table := newTestTable()
	for i := range 100 {
		table.Insert(testKey(i))
	}

	if table.Contains("jamais-insérée") {
		t.Error("une clé jamais insérée ne doit pas être trouvée")
	}
}

func TestTableReachesItsLoadLimitOnTheInsertThatExceedsSevenKeysPerGroup(t *testing.T) {
	table := newTestTable()

	for i := range 7 {
		table.Insert(testKey(i))
	}

	if got := table.loadLimit(); got != 7 {
		t.Errorf("loadLimit = %d, voulu 7 pour un groupe", got)
	}
	if !table.reachesLoadLimit() {
		t.Error("le 8e ajout dépasse 7 clés dans un groupe, la table doit grossir")
	}
}

func TestTableIsAtMaxCapacityOnlyOnceItHoldsMaxTableCapacitySlots(t *testing.T) {
	table := newTestTable()
	if table.atMaxCapacity() {
		t.Error("une table d'un groupe n'est pas au plafond")
	}

	for i := range maxKeysPerTable {
		table.Insert(testKey(i))
	}

	if !table.atMaxCapacity() {
		t.Errorf("une table de %d slots est au plafond", table.Capacity())
	}
}

func TestProbeVisitsEveryGroupOnceStartingFromTheGroupOfTheHash(t *testing.T) {
	table := newTable(seed, 4, 0)
	const hash = 2 << 7

	var visited []*Group
	for g := range table.probe(hash) {
		visited = append(visited, g)
	}

	want := []*Group{&table.groups[2], &table.groups[3], &table.groups[0], &table.groups[1]}
	if !slices.Equal(visited, want) {
		t.Errorf("ordre de sondage inattendu : %v, voulu %v", visited, want)
	}
}

func TestProbeStopsAsSoonAsTheCallerStopsIterating(t *testing.T) {
	table := newTable(seed, 4, 0)

	visits := 0
	for range table.probe(0) {
		visits++
		break
	}

	if visits != 1 {
		t.Errorf("%d groupes visités après un break, voulu 1", visits)
	}
}
