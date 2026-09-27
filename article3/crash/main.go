// Command crash lance deux goroutines qui écrivent chacune 1000 fois dans
// la même map, sans synchronisation. Le runtime détecte l'accès concurrent
// et termine tout le programme avec `fatal error: concurrent map writes`.
// Ce fatal error n'est pas une panic : même un recover posé directement
// dans la goroutine qui écrit ne l'intercepte jamais, ce que ce programme
// démontre en s'exécutant jusqu'au crash.
package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[int]int)

	var wg sync.WaitGroup
	wg.Add(2)

	writer := func(start int) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover a intercepté:", r)
			}
		}()
		for i := range 1000 {
			m[start+i] = i
		}
	}

	go writer(0)
	go writer(1000)

	wg.Wait()
	fmt.Println("terminé sans crash, taille finale:", len(m))
}
