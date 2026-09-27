// Benchmarks comparant Mutex, RWMutex et sync.Map sous différents ratios
// lecture/écriture. Les chiffres dépendent de la version de Go et de la
// machine : relancez-les vous-même plutôt que de vous fier à des chiffres
// figés dans un article. Voir le README pour la commande.
package article3

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
)

const benchKeyCount = 100

type mutexMap struct {
	mu sync.Mutex
	m  map[int]int
}

func newMutexMap() *mutexMap {
	return &mutexMap{m: make(map[int]int, benchKeyCount)}
}

func (s *mutexMap) read(key int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key]
}

func (s *mutexMap) write(key, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

type rwMutexMap struct {
	mu sync.RWMutex
	m  map[int]int
}

func newRWMutexMap() *rwMutexMap {
	return &rwMutexMap{m: make(map[int]int, benchKeyCount)}
}

func (s *rwMutexMap) read(key int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[key]
}

func (s *rwMutexMap) write(key, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func newBenchSyncMap() *sync.Map {
	m := &sync.Map{}
	for i := range benchKeyCount {
		m.Store(i, i)
	}
	return m
}

var writePercentages = []int{0, 10, 50, 90}

func BenchmarkMutexMap(b *testing.B) {
	for _, writeRatioPercent := range writePercentages {
		b.Run(fmt.Sprintf("write%d%%", writeRatioPercent), func(b *testing.B) {
			s := newMutexMap()
			for i := range benchKeyCount {
				s.write(i, i)
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := rand.IntN(benchKeyCount)
					if rand.IntN(100) < writeRatioPercent {
						s.write(key, i)
					} else {
						s.read(key)
					}
					i++
				}
			})
		})
	}
}

func BenchmarkRWMutexMap(b *testing.B) {
	for _, writeRatioPercent := range writePercentages {
		b.Run(fmt.Sprintf("write%d%%", writeRatioPercent), func(b *testing.B) {
			s := newRWMutexMap()
			for i := range benchKeyCount {
				s.write(i, i)
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := rand.IntN(benchKeyCount)
					if rand.IntN(100) < writeRatioPercent {
						s.write(key, i)
					} else {
						s.read(key)
					}
					i++
				}
			})
		})
	}
}

func BenchmarkSyncMap(b *testing.B) {
	for _, writeRatioPercent := range writePercentages {
		b.Run(fmt.Sprintf("write%d%%", writeRatioPercent), func(b *testing.B) {
			s := newBenchSyncMap()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := rand.IntN(benchKeyCount)
					if rand.IntN(100) < writeRatioPercent {
						s.Store(key, i)
					} else {
						s.Load(key)
					}
					i++
				}
			})
		})
	}
}
