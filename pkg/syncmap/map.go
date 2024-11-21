package syncmap

import "sync"

// SyncMap is a thread-safe map.
type SyncMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSyncMap creates a new SyncMap.
func NewSyncMap[K comparable, V any](n int) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: make(map[K]V, n),
	}
}

// Load loads a value from the map.
func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.m[key]
	return val, ok
}

// Store stores a value in the map.
func (s *SyncMap[K, V]) Store(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}

// Has returns true if the key exists in the map.
func (s *SyncMap[K, V]) Has(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.m[key]
	return ok
}

// Delete deletes a value from the map.
func (s *SyncMap[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, key)
}

// Range calls f sequentially for each key and value present in the map.
func (s *SyncMap[K, V]) Range(callback func(key K, value V) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.m {
		if !callback(k, v) {
			break
		}
	}
}

// Len returns the length of the map.
func (s *SyncMap[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// ToMap returns a copy of the map.
func (s *SyncMap[K, V]) ToMap() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyMap := make(map[K]V, len(s.m))
	for k, v := range s.m {
		copyMap[k] = v
	}
	return copyMap
}
