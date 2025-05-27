package Transposition

import (
	"container/list"
	"sync"
)

// TranspositionTable implements an OrderedDict-like structure with LRU eviction
type TranspositionTable struct {
	maxSize int
	items   map[uint64]*list.Element
	evictList *list.List
	mutex   sync.RWMutex // For thread safety
}

// entry is used to store the key-value pair in the eviction list
type entry struct {
	key   uint64
	value interface{}
}

// NewTranspositionTable creates a new TranspositionTable with the given max size
func NewTranspositionTable(maxSize int) *TranspositionTable {
	return &TranspositionTable{
		maxSize:   maxSize,
		items:     make(map[uint64]*list.Element),
		evictList: list.New(),
	}
}

// Get retrieves a value from the table and moves it to the front (most recently used)
func (t *TranspositionTable) Get(key uint64) (interface{}, bool) {
	t.mutex.RLock()
	element, exists := t.items[key]
	t.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	// Move to front (mark as most recently used)
	t.mutex.Lock()
	t.evictList.MoveToFront(element)
	t.mutex.Unlock()

	return element.Value.(*entry).value, true
}

// Put adds or updates a value in the table
func (t *TranspositionTable) Put(key uint64, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// If key exists, update its value and move to front
	if element, exists := t.items[key]; exists {
		t.evictList.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}

	// Add new item
	element := t.evictList.PushFront(&entry{key: key, value: value})
	t.items[key] = element

	// Evict oldest if we're over capacity
	if t.evictList.Len() > t.maxSize {
		t.removeOldest()
	}
}

// Contains checks if a key exists in the table
func (t *TranspositionTable) Contains(key uint64) bool {
	t.mutex.RLock()
	_, exists := t.items[key]
	t.mutex.RUnlock()
	return exists
}

// Len returns the current size of the table
func (t *TranspositionTable) Len() int {
	t.mutex.RLock()
	length := len(t.items)
	t.mutex.RUnlock()
	return length
}

// removeOldest removes the oldest item from the cache
func (t *TranspositionTable) removeOldest() {
	oldest := t.evictList.Back()
	if oldest != nil {
		entry := oldest.Value.(*entry)
		delete(t.items, entry.key)
		t.evictList.Remove(oldest)
	}
}

// Clear empties the table
func (t *TranspositionTable) Clear() {
	t.mutex.Lock()
	t.items = make(map[uint64]*list.Element)
	t.evictList.Init()
	t.mutex.Unlock()
}