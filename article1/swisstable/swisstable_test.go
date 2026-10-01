package swisstable

import (
	"hash/maphash"
	"testing"
)

var seed = maphash.MakeSeed()

func hashKey(key string) uint64 {
	return maphash.String(seed, key)
}

func TestH1AndH2PartitionA64BitHash(t *testing.T) {
	const (
		hash   = uint64(0xABCD_1234_5678_90EF)
		wantH1 = uint64(0x1579a2468acf121)
		wantH2 = byte(0x6f)
	)

	gotH1 := h1(hash)
	gotH2 := h2(hash)

	if gotH1 != wantH1 {
		t.Errorf("h1(%x) = %x, voulu %x", hash, gotH1, wantH1)
	}
	if gotH2 != wantH2 {
		t.Errorf("h2(%x) = %x, voulu %x", hash, gotH2, wantH2)
	}

	if reconstructed := gotH1<<7 | uint64(gotH2); reconstructed != hash {
		t.Errorf("h1(%x)<<7 | h2(%x) = %x, voulu %x", hash, hash, reconstructed, hash)
	}
}

func TestLookupFindsInsertedKey(t *testing.T) {
	g := NewEmptyGroup()
	hash := hashKey("marie")

	if !g.Insert(hash, "marie") {
		t.Fatal("insertion refusée dans un groupe vide")
	}

	found, slot := g.Lookup(hash, "marie")
	if !found {
		t.Fatal("clé insérée non retrouvée")
	}
	if g.keys[slot] != "marie" {
		t.Errorf("slot %d contient %q, voulu marie", slot, g.keys[slot])
	}
}

func TestLookupRejectsAbsentKey(t *testing.T) {
	g := NewEmptyGroup()
	g.Insert(hashKey("marie"), "marie")

	found, _ := g.Lookup(hashKey("julien"), "julien")
	if found {
		t.Error("une clé jamais insérée ne doit pas être trouvée")
	}
}

func TestInsertFailsOnFullGroup(t *testing.T) {
	g := NewEmptyGroup()
	for i := range groupSize {
		key := string(rune('a' + i))
		if !g.Insert(hashKey(key), key) {
			t.Fatalf("insertion %d/%d a échoué avant que le groupe soit plein", i+1, groupSize)
		}
	}

	if g.Insert(hashKey("overflow"), "overflow") {
		t.Error("un groupe de 8 slots pleins ne doit plus accepter d'insertion")
	}
}

// withSameH2ButDifferentH1 fabrique un hash différent de hash qui conserve son H2 (les 7 bits de poids faible).
//
//   - `hash &^ 0x7f` (AND NOT) met les 7 bits bas à zéro, ce qui efface H2 ;
//   - `^ 0x1_0000_0000` inverse le bit 32 : le hash change, donc H1 aussi. Tout bit au-dessus du bit 6 conviendrait ;
//   - `| hash & 0x7f` remet le H2 d'origine dans les 7 bits bas libérés.
//
// MatchH2 ne comparant que H2, le slot de la clé d'origine est remonté à tort (faux positif) : seule la
// comparaison de clé dans Lookup doit l'écarter.
func withSameH2ButDifferentH1(hash uint64) uint64 {
	return (hash&^0x7f ^ 0x1_0000_0000) | hash&0x7f
}

func TestMatchH2CanReturnFalsePositiveConfirmedByKeyComparison(t *testing.T) {
	g := NewEmptyGroup()
	hash := hashKey("marie")
	g.Insert(hash, "marie")

	collidingHash := withSameH2ButDifferentH1(hash)

	matches := g.MatchH2(collidingHash)
	if len(matches) == 0 {
		t.Fatal("attendu au moins un slot avec le même H2")
	}

	found, _ := g.Lookup(collidingHash, "quelqu-un-d-autre")
	if found {
		t.Error("un H2 identique ne doit pas suffire sans comparaison de clé")
	}
}

func TestGroupForStaysWithinBounds(t *testing.T) {
	const groupCount = 16

	for _, key := range []string{"a", "b", "c", "marie", "julien"} {
		g := GroupFor(hashKey(key), groupCount)
		if g >= groupCount {
			t.Errorf("GroupFor(%q) = %d, hors bornes [0,%d)", key, g, groupCount)
		}
	}
}
