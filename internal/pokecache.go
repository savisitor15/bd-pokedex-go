package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	timeToLive    time.Duration
	cacheMux      *sync.Mutex
	cacheInternal map[string]CacheEntry
	tick          *time.Ticker
	destroy       chan int
}

func (c Cache) Add(key string, val []byte) {
	c.cacheMux.Lock()
	c.cacheInternal[key] = CacheEntry{val: val, createdAt: time.Now()}
	c.cacheMux.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.cacheMux.Lock()
	defer c.cacheMux.Unlock() // unlock on return
	if bOut, ok := c.cacheInternal[key]; ok {
		return bOut.val, true
	}
	return nil, false
}

func (c Cache) reaploop() {
	for {
		select {
		case <-c.destroy:
			return
		case <-c.tick.C:
			c.cacheMux.Lock()
			// loop
			for k, elm := range c.cacheInternal {
				if time.Since(elm.createdAt) > c.timeToLive {
					// Expired
					delete(c.cacheInternal, k)
				}
			}
			c.cacheMux.Unlock()
		}
	}
}

func (c Cache) Destroy() {
	// End the cache
	c.destroy <- 1
	// dump the map
	c.cacheMux.Lock()
	for k := range c.cacheInternal {
		delete(c.cacheInternal, k)
	}
	c.cacheMux.Unlock()
}

func NewCache(interval time.Duration) Cache {
	var nCache Cache = Cache{timeToLive: interval, cacheMux: &sync.Mutex{}, cacheInternal: make(map[string]CacheEntry), tick: time.NewTicker(interval)}
	// Launch the cache maintainer
	go nCache.reaploop()
	return nCache
}
