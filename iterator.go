package nncache

import "sync"

// EntryInfoIterator allows to iterate over entries in the cache
type EntryInfoIterator struct {
	mutex         sync.Mutex
	cache         *NNCache
	currentShard  int
	currentIndex  int
	elements      []uint32
	elementsCount int
	valid         bool
}
