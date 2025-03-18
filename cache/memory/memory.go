package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofreego/goutils/datastructure"
)

type object struct {
	key     string
	value   any
	expriry time.Time
}

func compare(a, b *object) int {
	if a.expriry.Before(b.expriry) {
		return -1
	}
	if a.expriry.After(b.expriry) {
		return 1
	}
	return 0
}

type Cache struct {
	cache   map[string]*object
	minHeap datastructure.Heap[*object]
}

func NewCache() *Cache {

	cache := &Cache{
		cache:   make(map[string]*object),
		minHeap: datastructure.NewHeap(datastructure.MinHeap, compare),
	}
	go cache.autoRemoveExpiredKeys()
	return cache
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	v, ok := c.cache[key]
	if !ok {
		return "", nil
	}
	return v.value.(string), nil
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
	c.minHeap.Insert(&v)
	return nil
}

func (c *Cache) SetWithTimeout(ctx context.Context, key string, value any, timeout time.Duration) error {
	v := object{
		key:     key,
		value:   value,
		expriry: time.Now().Add(timeout),
	}
	c.cache[key] = &v
	c.minHeap.Insert(&v)
	return nil
}

func (c *Cache) autoRemoveExpiredKeys() {
	for {
		if c.minHeap.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}
		obj, err := c.minHeap.Peek()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if obj.expriry.Before(time.Now()) {
			c.minHeap.Extract()
			delete(c.cache, obj.key)
		} else {
			time.Sleep(time.Second)
		}
	}
}
