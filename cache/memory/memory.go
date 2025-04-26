package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofreego/ds"
)

type object struct {
	key     string
	value   any
	expriry time.Time
}

func less(a, b *object) bool {
	return a.expriry.Before(b.expriry)
}

type Cache struct {
	cache   map[string]*object
	minHeap ds.Heap[*object]
}

func NewCache() *Cache {

	cache := &Cache{
		cache:   make(map[string]*object),
		minHeap: ds.NewHeap(ds.MinHeap, less),
	}
	go cache.autoRemoveExpiredKeys()
	return cache
}

func (c *Cache) GetV(ctx context.Context, key string, value any) error {
	v, ok := c.cache[key]
	if !ok {
		return nil
	}
	bytes, err := json.Marshal(v.value)
	if err != nil {
		return fmt.Errorf("value is not compatible with given object,Err: %s", err.Error())
	}
	err = json.Unmarshal(bytes, value)
	if err != nil {
		return fmt.Errorf("value is not compatible with given object,Err: %s", err.Error())
	}
	return nil
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {

	v := object{
		key:   key,
		value: value,
		// 365*24*10 = 87600 = 10 years
		expriry: time.Now().Add(time.Hour * 87600),
	}
	c.cache[key] = &v
	c.minHeap.Push(&v)
	return nil
}

func (c *Cache) SetWithTimeout(ctx context.Context, key string, value any, timeout time.Duration) error {
	v := object{
		key:     key,
		value:   value,
		expriry: time.Now().Add(timeout),
	}
	c.cache[key] = &v
	c.minHeap.Push(&v)
	return nil
}

func (c *Cache) autoRemoveExpiredKeys() {
	for {
		if c.minHeap.Size() == 0 {
			time.Sleep(time.Second)
			continue
		}
		obj := c.minHeap.Top()
		if obj.expriry.Before(time.Now()) {
			c.minHeap.Pop()
			delete(c.cache, obj.key)
		} else {
			time.Sleep(time.Second)
		}
	}
}
