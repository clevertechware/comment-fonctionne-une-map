# Exemples Go — série « Comment fonctionne une map »

Dépôt d'exemples accompagnant les trois articles de la série. Un module Go,
un sous-dossier par article.

- `article1/` : [Ce qui se passe quand on écrit m[k] en Go](https://www.clevertechware.fr/blog/posts/2026/comment-fonctionne-une-map-en-go)
- `article2/` : [Quand la map grossit](https://www.clevertechware.fr/blog/posts/2026/quand-une-map-go-grossit)
- `article3/` : [Maps et goroutines : pourquoi ça crashe et que choisir](https://www.clevertechware.fr/blog/posts/2026/maps-go-et-goroutines)

## Prérequis

- Go 1.23+ (`go version`). Ce module cible `go 1.23` dans `go.mod` pour
  pouvoir comparer hmap et Swiss Tables avec deux toolchains différentes
  (voir la section `article2/` ci-dessous). `GOTOOLCHAIN=local` force
  la toolchain installée, même si `GOTOOLCHAIN` est défini ailleurs.

## Build, run, test

```bash
go build ./...
go vet ./...
go test ./...
go test ./... -race
```

## Contenu par article

### `article1/` — Ce qui se passe quand on écrit `m[k]`

- `lookup_bench_test.go` : coût d'un lookup réussi vs manqué sur
  `map[string]int`. Lance avec `-benchmem` pour voir les allocations :
  ```bash
  go test ./article1/... -bench=Lookup -benchmem -run=^$
  ```
- `lookup_test.go` : vérifie que ces deux lookups n'allouent rien
  (`testing.AllocsPerRun`).
- `swisstable/` : **réimplémentation pédagogique minimale**, pas le code du
  runtime. Le vrai code vit dans `internal/runtime/maps`, un paquet interne
  à la stdlib non importable depuis l'extérieur. Les constantes et formules
  (`h1`, `h2`, `ctrlEmpty`, `ctrlDeleted`) sont recopiées à l'identique
  depuis le runtime installé (Go 1.27.1, `go env GOROOT`, fichiers
  `src/internal/runtime/maps/map.go` et `group.go`) ; le reste (structure
  `Group`, `Insert`, `Lookup`) est une simplification pour l'article, testée
  dans `swisstable_test.go`.

### `article2/` — Quand la map grossit

- `alloc_test.go` : allocations d'une insertion de 1000 clés, avec et sans
  `make(map[K]V, 1000)`. Le nombre de buckets/tables n'est pas observable
  proprement depuis l'extérieur du runtime ; les allocations mémoire le
  sont, via `testing.AllocsPerRun` et `-benchmem`.
  ```bash
  go test ./article2/... -bench=Insert -benchmem -run=^$
  ```
- `insert_bench_test.go` : `BenchmarkInsert1000`, le bench à relancer sous
  deux toolchains pour comparer hmap (Go 1.23) et Swiss Tables (Go 1.24+ par
  défaut, seule implémentation depuis Go 1.26). Une seule toolchain
  (go1.27.1) est installée ici : aucun chiffre de comparaison n'est fourni
  dans l'article, volontairement. Pour comparer vous-même, forcez
  `GOTOOLCHAIN=local` sur les deux commandes : si `GOTOOLCHAIN` est défini
  dans l'environnement ou via `go env -w` (par exemple `go1.27.1`), Go
  basculerait sur cette toolchain et mesurerait les Swiss Tables deux fois.
  ```bash
  go install golang.org/dl/go1.23.12@latest
  go1.23.12 download
  GOTOOLCHAIN=local go1.23.12 test ./article2/... -bench=Insert1000 -benchtime=1s -run=^$   # go1.23.12, hmap

  GOTOOLCHAIN=local go test ./article2/... -bench=Insert1000 -benchtime=1s -run=^$   # go1.27.1, Swiss Tables
  ```
  Alternative sans changer de toolchain, sur Go 1.24/1.25 uniquement :
  `GOEXPERIMENT=noswissmap go test ...` bascule sur l'ancienne implémentation
  hmap. Ce flag n'existe plus en Go 1.26+.
- `iteration_test.go`, `range_example_test.go` : itération non déterministe
  (`TestTwoRangesOverSameMapCanProduceDifferentOrders`) et ordre stable via
  `slices.Sorted(maps.Keys(m))`.

### `article3/` — Maps et goroutines

- `crash/` : programme exécutable séparé (`main.go`) qui lance deux
  goroutines écrivant chacune 1000 fois dans la même map, provoquant
  `fatal error: concurrent map writes`. Le `recover` est posé directement
  dans la goroutine `writer` qui écrit, et ne l'intercepte quand même pas :
  ce fatal error n'est pas une panic récupérable.
  ```bash
  go run ./article3/crash    # crashe volontairement (le binaire sort en code 2 ; `go run` lui-même sort en 1)
  ```
  `main_test.go` lance ce programme en sous-processus et vérifie le message
  sur stderr, sans faire planter la suite de tests. Sous `-race`, ce test
  particulier est sauté : le timing change, et dans nos essais (5/5
  exécutions) le race detector signale une seule `WARNING: DATA RACE` sans
  que le fatal error du runtime se déclenche ; le programme va au bout et
  sort avec le code 66. Ce test vise spécifiquement le fatal error, donc il
  faut le lancer sans `-race`.
- `safecounter.go`, `synccache.go` : les extraits `SafeCounter` (Mutex) et
  `sync.Map` cités dans l'article, testés dans `safecounter_test.go` et
  `synccache_test.go`.
- `bench_test.go` : benchmark Mutex / RWMutex / `sync.Map` sur plusieurs
  ratios lecture/écriture (0 %, 10 %, 50 %, 90 % d'écritures).
  ```bash
  go test ./article3/... -bench=. -benchmem -run=^$
  ```
  Les chiffres dépendent de la version de Go et de la machine : relancez ce
  benchmark vous-même plutôt que de vous fier à des chiffres publiés
  dans un article.

## Vérification rapide (CI-like)

```bash
go build ./...
go vet ./...
go test ./...
go test ./... -race
go test ./... -bench=. -benchtime=100ms -run=^$ -benchmem
```
