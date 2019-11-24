package nncache

// NNCache is fast, concurrent, evicting cache created to keep big number of entries without impact on performance.
// It keeps entries on heap but omits GC for them. To achieve that, operations take place on byte arrays,
// therefore entries (de)serialization in front of the cache will be needed in most use cases.
type NNCache struct {
	shards       []*cacheShard
	hash         Hasher
	lifeWindow   uint64
	clock        clock
	config       Config
	shardMask    uint64
	maxShardSize uint32
	close        chan struct{}
}

const (
	// Expired means the key is past its LifeWindow.
	// @TODO: Go defaults to 0 so in case we want to return EntryStatus back to the caller Expired cannot be differentiated
	Expired RemoveReason = iota
	// NoSpace means the key is the oldest and the cache size was at its maximum when Set was called, or the
	// entry exceeded the maximum shard size.
	NoSpace
	// Deleted means Delete was called and this key was removed as a result.
	Deleted
)
const (
	minimumEntriesInShard = 10 // Minimum number of entries in single shard
)

// RemoveReason is a value used to signal to the user why a particular key was removed in the OnRemove callback.
type RemoveReason uint32

// NewNNCache initialize new instance of NNCache
func NewNNCache(config Config) (*NNCache, error) {
	return newNNCache(config, &systemClock{})
}

func newNNCache(config Config, clock clock) (*NNCache, error) {
	if config.Hasher == nil {
		config.Hasher = newDefaultHasher()
	}
	cache := &NNCache{
		shards:       make([]*cacheShard, config.Shards),
		hash:         config.Hasher,
		lifeWindow:   uint64(config.LifeWindow.Seconds()),
		clock:        clock,
		config:       config,
		shardMask:    uint64(config.Shards - 1),
		maxShardSize: uint32(config.maximumShardSize()),
		close:        make(chan struct{}),
	}
	for i := 0; i < config.Shards; i++ {
		cache.shards[i] = initNewShard(config, nil, clock)
	}
	return cache, nil
}

// Close is used to signal a shutdown of the cache when you are done with it.
// This allows the cleaning goroutines to exit and ensures references are not
// kept to the cache preventing GC of the entire cache.
func (c *NNCache) Close() error {
	close(c.close)
	return nil
}

// Set saves entry under the key
func (c *NNCache) Set(key string, entry []byte) error {
	hashedKey := c.hash.Sum64(key)
	shard := c.getShard(hashedKey)
	return shard.set(key, hashedKey, entry)
}

// Get reads entry for the key.
// It returns an ErrEntryNotFound when
// no entry exists for the given key.
func (c *NNCache) Get(key string) ([]byte, error) {
	hashedKey := c.hash.Sum64(key)
	shard := c.getShard(hashedKey)
	return shard.get(key, hashedKey)
}

// Delete removes the key
func (c *NNCache) Delete(key string) error {
	return nil
}

// Reset empties all cache shards
func (c *NNCache) Reset() error {

	return nil
}

// Len computes number of entries in cache
func (c *NNCache) Len() int {
	return 0
}

// Capacity returns amount of bytes store in the cache.
func (c *NNCache) Capacity() int {
	return 0
}

// Stats returns cache's statistics
func (c *NNCache) Stats() Stats {
	var s Stats

	return s
}

// Iterator returns iterator function to iterate over EntryInfo's from whole cache.
func (c *NNCache) Iterator() *EntryInfoIterator {
	return nil
}

func (c *NNCache) getShard(hashedKey uint64) (shard *cacheShard) {
	return c.shards[hashedKey&c.shardMask]
}
