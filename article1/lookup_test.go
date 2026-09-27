package article1

import "testing"

func TestLookupHitAndMissAllocateNothing(t *testing.T) {
	m := buildStringIntMap(1000)

	hitAllocs := testing.AllocsPerRun(100, func() {
		_ = m["key-0500"]
	})
	if hitAllocs != 0 {
		t.Errorf("lookup réussi: %.1f allocs/op, voulu 0", hitAllocs)
	}

	missAllocs := testing.AllocsPerRun(100, func() {
		_ = m["key-9999"]
	})
	if missAllocs != 0 {
		t.Errorf("lookup manqué: %.1f allocs/op, voulu 0", missAllocs)
	}
}
