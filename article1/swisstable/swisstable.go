// Package swisstable est une réimplémentation pédagogique minimale du
// filtre H1/H2 des Swiss Tables de Go. Ce n'est PAS le code du runtime :
// c'est une version simplifiée, pour illustrer le découpage du hash et le
// test du mot de contrôle. table.go et map.go y ajoutent la croissance décrite
// dans l'article 2 (doublement sur place d'une table, scission à 1024 slots,
// annuaire). Le code réel vit dans internal/runtime/maps
// (non importable depuis l'extérieur de la stdlib) ; les formules h1/h2
// et les constantes ctrlEmpty/ctrlDeleted ci-dessous sont recopiées à
// l'identique depuis internal/runtime/maps/map.go et group.go, Go 1.27.1
// (go env GOROOT -> src/internal/runtime/maps).
package swisstable

const groupSize = 8

const (
	ctrlEmpty   byte = 0b10000000
	ctrlDeleted byte = 0b11111110
)

// h1 retourne les 57 bits hauts du hash : le numéro de groupe.
func h1(h uint64) uint64 {
	return h >> 7
}

// h2 retourne les 7 bits bas du hash : la valeur stockée dans l'octet de contrôle d'un slot occupé.
func h2(h uint64) byte {
	return byte(h & 0x7f)
}

// Group est un groupe de 8 slots avec son mot de contrôle associé.
type Group struct {
	ctrl [groupSize]byte
	keys [groupSize]string
}

// NewEmptyGroup construit un groupe dont tous les slots sont vides.
func NewEmptyGroup() *Group {
	g := &Group{}
	g.markAllEmpty()
	return g
}

func (g *Group) markAllEmpty() {
	for i := range g.ctrl {
		g.ctrl[i] = ctrlEmpty
	}
}

// Insert place key dans le premier slot vide du groupe, en y posant son H2.
// Retourne false si le groupe est plein.
func (g *Group) Insert(hash uint64, key string) bool {
	for i, c := range g.ctrl {
		if c == ctrlEmpty || c == ctrlDeleted {
			g.ctrl[i] = h2(hash)
			g.keys[i] = key
			return true
		}
	}
	return false
}

// MatchH2 retourne les indices des slots dont l'octet de contrôle égal H2.
//
// Comme dans le runtime, une correspondance H2 est un candidat : elle doit encore être confirmée par une comparaison
// de la clé complète, parce que H2 ne fait que 7 bits et peut donner de faux positifs.
func (g *Group) MatchH2(hash uint64) []int {
	target := h2(hash)
	var matches []int
	for i, c := range g.ctrl {
		if c == target {
			matches = append(matches, i)
		}
	}
	return matches
}

// Lookup applique le filtre H2 puis confirme sur la clé complète, exactement
// comme le fait matchH2 suivi de la comparaison de clé dans le runtime.
func (g *Group) Lookup(hash uint64, key string) (found bool, slot int) {
	for _, i := range g.MatchH2(hash) {
		if g.keys[i] == key {
			return true, i
		}
	}
	return false, -1
}

// GroupFor calcule H1(hash) pour choisir le numéro de groupe dans une table
// de groupCount groupes (groupCount doit être une puissance de 2).
func GroupFor(hash uint64, groupCount uint64) uint64 {
	return h1(hash) & (groupCount - 1)
}
