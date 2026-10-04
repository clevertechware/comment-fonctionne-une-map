package swisstable

import "testing"

func fillMap(m *Map, n int) {
	for i := range n {
		m.Insert(testKey(i))
	}
}

func TestTableAtMaxCapacityRefusesInsertInsteadOfGrowing(t *testing.T) {
	table := newTestTable()
	for i := range maxKeysPerTable {
		if !table.Insert(testKey(i)) {
			t.Fatalf("insertion %d refusée avant le plafond", i+1)
		}
	}

	if table.Insert(testKey(maxKeysPerTable)) {
		t.Error("une table de 1024 slots à 7/8 de charge doit refuser l'insertion pour se faire scinder")
	}
	if got := table.Capacity(); got != maxTableCapacity {
		t.Errorf("capacité = %d, voulu %d : une table ne grossit pas au-delà du plafond", got, maxTableCapacity)
	}
}

func TestMapSplitsItsTableOnceTheFirstOneIsFull(t *testing.T) {
	m := NewMap()

	fillMap(m, maxKeysPerTable)
	if got := m.TableCount(); got != 1 {
		t.Fatalf("%d tables pour %d clés, voulu 1", got, maxKeysPerTable)
	}

	m.Insert(testKey(maxKeysPerTable))
	if got := m.TableCount(); got != 2 {
		t.Errorf("%d tables après une clé de plus que le plafond d'une table, voulu 2", got)
	}
	if got := len(m.directory); got != 2 {
		t.Errorf("annuaire de %d entrées après la première scission, voulu 2", got)
	}
}

func TestMapKeepsEveryKeyAcrossSplits(t *testing.T) {
	const n = 20_000
	m := NewMap()

	fillMap(m, n)

	for i := range n {
		if !m.Contains(testKey(i)) {
			t.Fatalf("clé %q perdue pendant les scissions", testKey(i))
		}
	}
	if got := m.Len(); got != n {
		t.Errorf("Len = %d, voulu %d", got, n)
	}
}

func TestMapRejectsAbsentKeyAfterSplits(t *testing.T) {
	m := NewMap()
	fillMap(m, 5000)

	if m.Contains("jamais-insérée") {
		t.Error("une clé jamais insérée ne doit pas être trouvée")
	}
}

func TestSplitReplacesOnlyTheFullTableAndLeavesTheOthersUntouched(t *testing.T) {
	m := NewMap()
	fillMap(m, 5000)
	before := m.distinctTables()

	next := 5000
	for m.TableCount() == len(before) {
		m.Insert(testKey(next))
		next++
	}

	after := map[*Table]bool{}
	for _, table := range m.distinctTables() {
		after[table] = true
	}
	replaced := 0
	for _, table := range before {
		if !after[table] {
			replaced++
		}
	}
	if replaced != 1 {
		t.Errorf("%d tables remplacées par la scission, voulu exactement 1", replaced)
	}
}

func TestDirectoryDoublesOnlyWhenTheSplitTableReadsAsManyBitsAsTheDirectory(t *testing.T) {
	m := NewMap()
	fillMap(m, 20_000)

	if got, want := len(m.directory), 1<<m.globalDepth; got != want {
		t.Fatalf("annuaire de %d entrées pour globalDepth=%d, voulu %d", got, m.globalDepth, want)
	}
	for _, table := range m.distinctTables() {
		if table.localDepth > m.globalDepth {
			t.Errorf("localDepth=%d dépasse globalDepth=%d", table.localDepth, m.globalDepth)
		}
	}
}

func TestTopBitsReadsTheHighestBitsOfTheHash(t *testing.T) {
	const hash = 0b1011 << 60

	tests := []struct {
		count uint8
		want  uint64
	}{
		{count: 0, want: 0},
		{count: 1, want: 0b1},
		{count: 4, want: 0b1011},
		{count: 64, want: hash},
	}
	for _, tt := range tests {
		if got := topBits(hash, tt.count); got != tt.want {
			t.Errorf("topBits(%#x, %d) = %#b, voulu %#b", uint64(hash), tt.count, got, tt.want)
		}
	}
}
