package nncache

import "time"

// Config for NNCache
type Config struct {
	// Hasher used to map between string keys and unsigned 64bit integers, by default fnv64 hashing is used.
	Hasher Hasher
	// Number of cache shards, value must be a power of two
	Shards int
	// Time after which entry can be evicted
	LifeWindow time.Duration
	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed. Setting to < 1 second is counterproductive — nncache has a one second resolution.
	CleanWindow time.Duration
	// Max number of entries in life window. Used only to calculate initial size for cache shards.
	// When proper value is set then additional memory allocation does not occur.
	MaxEntriesInWindow int
	// Max size of entry in bytes. Used only to calculate initial size for cache shards.
	MaxEntrySize int
	// StatsEnabled if true calculate the number of times a cached resource was requested.
	StatsEnabled bool
	// Verbose mode prints information about new memory allocation
	Verbose bool
	// HardMaxCacheSize is a limit for cache size in MB. Cache will not allocate more memory than this limit.
	// It can protect application from consuming all available memory on machine, therefore from running OOM Killer.
	// Default value is 0 which means unlimited size. When the limit is higher than 0 and reached then
	// the oldest entries are overridden for the new ones.
	HardMaxCacheSize int
	// OnRemove is a callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	// ignored if OnRemoveWithMetadata is specified.
	OnRemove       func(key string, entry []byte)
	onRemoveFilter int
}

// DefaultConfig initializes config with default values.
// When load for NNCache can be predicted in advance then it is better to use custom config.
func DefaultConfig(eviction time.Duration) Config {
	return Config{
		Hasher:             newDefaultHasher(),
		Shards:             1024,
		LifeWindow:         eviction,
		CleanWindow:        1 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		StatsEnabled:       false,
		Verbose:            true,
		HardMaxCacheSize:   0,
	}
}

// maximumShardSize computes maximum shard size
func (c Config) maximumShardSize() int {
	maxShardSize := 0

	if c.HardMaxCacheSize > 0 {
		maxShardSize = convertMBToBytes(c.HardMaxCacheSize) / c.Shards
	}

	return maxShardSize
}
// initialShardSize computes initial shard size
func (c Config) initialShardSize() int {
	return max(c.MaxEntriesInWindow/c.Shards, minimumEntriesInShard)
}