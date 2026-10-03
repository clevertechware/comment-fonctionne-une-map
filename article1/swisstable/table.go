package swisstable

import "hash/maphash"

// maxAvgGroupLoad est le nombre maximal de slots occupés par groupe de 8, soit un facteur de charge de 7/8,
// recopié de internal/runtime/maps/map.go (Go 1.27.1).
const maxAvgGroupLoad = 7

// Table est un tableau de groupes qui double sur place quand son facteur de charge dépasse 7/8.
type Table struct {
	seed   maphash.Seed
	groups []Group
	used   int
}

// NewTable construit une table d'un seul groupe.
func NewTable(seed maphash.Seed) *Table {
	return &Table{seed: seed, groups: newGroups(1)}
}

func newGroups(count int) []Group {
	groups := make([]Group, count)
	for i := range groups {
		groups[i] = *NewEmptyGroup()
	}
	return groups
}

// Capacity retourne le nombre de slots de la table.
func (t *Table) Capacity() int {
	return len(t.groups) * groupSize
}

// Len retourne le nombre de clés de la table.
func (t *Table) Len() int {
	return t.used
}

func (t *Table) hash(key string) uint64 {
	return maphash.String(t.seed, key)
}

// Contains sonde les groupes à partir de GroupFor et s'arrête au premier groupe qui a encore un slot vide :
// une clé insérée aurait été placée avant lui.
func (t *Table) Contains(key string) bool {
	hash := t.hash(key)
	start := GroupFor(hash, uint64(len(t.groups)))
	for probe := range uint64(len(t.groups)) {
		g := &t.groups[(start+probe)%uint64(len(t.groups))]
		if found, _ := g.Lookup(hash, key); found {
			return true
		}
		if g.hasEmptySlot() {
			return false
		}
	}
	return false
}

func (g *Group) hasEmptySlot() bool {
	for _, c := range g.ctrl {
		if c == ctrlEmpty {
			return true
		}
	}
	return false
}

// Insert ajoute key à la table, en doublant sa taille si l'ajout dépasse le facteur de charge de 7/8.
func (t *Table) Insert(key string) {
	if t.Contains(key) {
		return
	}
	if t.used+1 > t.Capacity()*maxAvgGroupLoad/groupSize {
		t.grow()
	}
	t.place(t.hash(key), key)
	t.used++
}

func (t *Table) place(hash uint64, key string) {
	start := GroupFor(hash, uint64(len(t.groups)))
	for probe := uint64(0); ; probe++ {
		if t.groups[(start+probe)%uint64(len(t.groups))].Insert(hash, key) {
			return
		}
	}
}

// grow double le nombre de groupes et replace chaque clé : son groupe de départ dépend du nombre de groupes.
func (t *Table) grow() {
	old := t.groups
	t.groups = newGroups(2 * len(old))
	eachKey(old, func(key string) {
		t.place(t.hash(key), key)
	})
}

func eachKey(groups []Group, fn func(key string)) {
	for i := range groups {
		for slot, c := range groups[i].ctrl {
			if c != ctrlEmpty && c != ctrlDeleted {
				fn(groups[i].keys[slot])
			}
		}
	}
}
