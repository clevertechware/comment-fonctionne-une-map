package swisstable

import "hash/maphash"

// Map route chaque clé vers une Table grâce à un annuaire indexé par les premiers bits de son hash.
// Quand une table atteint son plafond, elle seule se scinde en deux : les autres ne bougent pas.
type Map struct {
	seed        maphash.Seed
	directory   []*Table
	globalDepth uint8
}

// NewMap construit une map d'une seule table.
func NewMap() *Map {
	seed := maphash.MakeSeed()
	return &Map{seed: seed, directory: []*Table{NewTable(seed)}}
}

// directoryIndex lit les globalDepth premiers bits du hash. Un décalage de 64 donne 0 en Go,
// ce qui couvre le cas globalDepth = 0.
func (m *Map) directoryIndex(hash uint64) uint64 {
	return hash >> (64 - m.globalDepth)
}

func (m *Map) tableFor(key string) *Table {
	return m.directory[m.directoryIndex(maphash.String(m.seed, key))]
}

// Contains retourne true si key a été insérée.
func (m *Map) Contains(key string) bool {
	return m.tableFor(key).Contains(key)
}

// Insert ajoute key, en scindant sa table autant de fois que nécessaire jusqu'à ce qu'elle l'accepte.
func (m *Map) Insert(key string) {
	for !m.tableFor(key).Insert(key) {
		m.split(m.tableFor(key))
	}
}

// Len retourne le nombre total de clés.
func (m *Map) Len() int {
	total := 0
	for _, t := range m.distinctTables() {
		total += t.Len()
	}
	return total
}

// TableCount retourne le nombre de tables distinctes : plusieurs entrées de l'annuaire peuvent viser la même.
func (m *Map) TableCount() int {
	return len(m.distinctTables())
}

func (m *Map) distinctTables() []*Table {
	var tables []*Table
	seen := map[*Table]bool{}
	for _, t := range m.directory {
		if !seen[t] {
			seen[t] = true
			tables = append(tables, t)
		}
	}
	return tables
}

// split remplace full par deux tables de même capacité, une par valeur du bit de hash qui suit ceux déjà lus
// par full. L'annuaire double d'abord si full lit déjà tous les bits qu'il lit lui-même.
func (m *Map) split(full *Table) {
	depth := full.localDepth + 1
	left := newTable(m.seed, len(full.groups), depth)
	right := newTable(m.seed, len(full.groups), depth)

	eachKey(full.groups, func(key string) {
		hash := maphash.String(m.seed, key)
		target := left
		if hash>>(64-depth)&1 == 1 {
			target = right
		}
		target.place(hash, key)
		target.used++
	})

	if depth > m.globalDepth {
		m.doubleDirectory()
	}
	for i, t := range m.directory {
		if t != full {
			continue
		}
		if i>>(m.globalDepth-depth)&1 == 0 {
			m.directory[i] = left
		} else {
			m.directory[i] = right
		}
	}
}

func (m *Map) doubleDirectory() {
	doubled := make([]*Table, 2*len(m.directory))
	for i, t := range m.directory {
		doubled[2*i] = t
		doubled[2*i+1] = t
	}
	m.directory = doubled
	m.globalDepth++
}
