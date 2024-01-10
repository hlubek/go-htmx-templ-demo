package service

import "sync"

type Counter interface {
	Increment() int
	Get() int
}

type InMemoryCounter struct {
	count int
	mx    sync.Mutex
}

func NewInMemoryCounter() *InMemoryCounter {
	return &InMemoryCounter{}
}

func (c *InMemoryCounter) Increment() int {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.count++
	return c.count
}

func (c *InMemoryCounter) Get() int {
	c.mx.Lock()
	defer c.mx.Unlock()

	return c.count
}
