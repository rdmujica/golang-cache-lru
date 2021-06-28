package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

const _max = 1000

type Cache struct {
	cache      *lru.Cache
	cacheMutex sync.Mutex
}

func NewCache(maxEntries int) *Cache {
	return &Cache{
		cache:      lru.New(maxEntries),
		cacheMutex: sync.Mutex{},
	}
}

func (c *Cache) Add(key interface{}, value interface{}) {
	c.cacheMutex.Lock()
	c.cache.Add(key, value)
	c.cacheMutex.Unlock()
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.cacheMutex.Lock()
	value, ok := c.cache.Get(key)
	c.cacheMutex.Unlock()
	return value, ok
}

func (c *Cache) Len() int {
	return c.cache.Len()
}

func (c *Cache) RemoveOldest() {
	c.cacheMutex.Lock()
	c.cache.RemoveOldest()
	c.cacheMutex.Unlock()
}

func (c *Cache) Clear() {
	c.cacheMutex.Lock()
	c.cache.Clear()
	c.cacheMutex.Unlock()
}

type producer struct {
	cache *Cache
}

func NewProducer(cahe *Cache) *producer {
	return &producer{
		cache: cahe,
	}
}

func (p *producer) Run() {
	for i := 0; i < _max; i++ {
		p.cache.Add(GetColumnName(i), i)
	}
}

type consumer struct {
	cache *Cache
}

func NewConsumer(cahe *Cache) *consumer {
	return &consumer{
		cache: cahe,
	}
}

func (p *consumer) Run(done chan bool) {
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < _max; i++ {
		col := GetColumnName(i)
		value, found := p.cache.Get(col)
		if !found {
			panic("missing value")
		}
		fmt.Println(col, value)
	}
	done <- true
}

func GetColumnName(i int) string {
	var col string
	for i > 0 {
		col = string(rune((i-1)%26+65)) + col
		i = (i - 1) / 26
	}
	return col
}

func main() {
	done := make(chan bool)
	cache := NewCache(0)
	prod := NewProducer(cache)
	cons := NewConsumer(cache)
	go prod.Run()
	go cons.Run(done)
	<-done
}
