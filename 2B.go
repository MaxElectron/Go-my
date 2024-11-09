//go:build !solution

package lrucache

import (
	"container/list"
)

type Entry struct {
	key   int
	value int
}

type LRUCache struct {
	capacity int
	cache    *list.List
	access   map[int]*list.Element
}

func New(capacity int) Cache {
	return &LRUCache{
		capacity: capacity,
		cache:    list.New(),
		access:   map[int]*list.Element{},
	}
}

func (c *LRUCache) Get(key int) (int, bool) {

	// Check for existence
	element, ok := c.access[key]
	if !ok {
		return 0, false
	}

	// Update the access time
	c.cache.MoveToBack(element)
	return element.Value.(Entry).value, true
}

func (c *LRUCache) Set(key, value int) {

	// If capacity is 0 do nothing
	if c.capacity == 0 {
		return
	}

	// If the key is present just update the value
	if element, ok := c.access[key]; ok {
		element.Value = Entry{key, value}
		c.cache.MoveToBack(element)
		return
	}

	// If the key is missing add the entry to the data and check the capacity
	element := c.cache.PushBack(Entry{key, value})
	c.access[key] = element
	if c.capacity < c.cache.Len() {
		delete(c.access, c.cache.Front().Value.(Entry).key)
		c.cache.Remove(c.cache.Front())
	}
}

func (c *LRUCache) Range(f func(key, value int) bool) {
	for element := c.cache.Front(); element != nil; element = element.Next() {
		if !f(element.Value.(Entry).key, element.Value.(Entry).value) {
			return
		}
	}
}

func (c *LRUCache) Clear() {
	c.cache.Init()
	c.access = map[int]*list.Element{}
}
