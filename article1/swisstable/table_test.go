package swisstable

import (
	"fmt"
	"hash/maphash"
	"testing"
)

func newTestTable() *Table {
	return NewTable(maphash.MakeSeed())
}

func key(i int) string {
	return fmt.Sprintf("cle-%d", i)
}

func TestTableDoublesItsCapacityWhenLoadExceedsSevenEighths(t *testing.T) {
	table := newTestTable()

	for i := range 7 {
		table.Insert(key(i))
	}
	if got := table.Capacity(); got != 8 {
		t.Fatalf("capacité après 7 clés = %d, voulu 8 (7/8 de charge n'est pas dépassé)", got)
	}

	table.Insert(key(7))
	if got := table.Capacity(); got != 16 {
		t.Errorf("capacité après 8 clés = %d, voulu 16", got)
	}
}

func TestTableKeepsEveryKeyAfterGrowing(t *testing.T) {
	const n = 1000
	table := newTestTable()

	for i := range n {
		table.Insert(key(i))
	}

	for i := range n {
		if !table.Contains(key(i)) {
			t.Fatalf("clé %q perdue pendant la croissance", key(i))
		}
	}
	if got := table.Len(); got != n {
		t.Errorf("Len = %d, voulu %d", got, n)
	}
}

func TestTableNeverExceedsSevenEighthsLoadAfterAnInsert(t *testing.T) {
	table := newTestTable()

	for i := range 1000 {
		table.Insert(key(i))
		if table.Len()*groupSize > table.Capacity()*maxAvgGroupLoad {
			t.Fatalf("après %d clés, %d/%d slots occupés dépasse 7/8", i+1, table.Len(), table.Capacity())
		}
	}
}

func TestInsertingAnExistingKeyChangesNeitherLenNorCapacity(t *testing.T) {
	table := newTestTable()
	for i := range 7 {
		table.Insert(key(i))
	}

	table.Insert(key(0))

	if table.Len() != 7 || table.Capacity() != 8 {
		t.Errorf("Len=%d Capacity=%d après réinsertion, voulu 7 et 8", table.Len(), table.Capacity())
	}
}

func TestContainsRejectsAbsentKeyAfterGrowing(t *testing.T) {
	table := newTestTable()
	for i := range 100 {
		table.Insert(key(i))
	}

	if table.Contains("jamais-insérée") {
		t.Error("une clé jamais insérée ne doit pas être trouvée")
	}
}
