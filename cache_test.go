package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c, _ := NewCache(100)
	for i := 0; i < 200; i++ {
		c.Set(i, i*i)
	}
	for i := 0; i < 100; i++ {
		res, ok := c.Get(i)
		if ok {
			t.Error(i, res, ok)
		}
	}
	for i := 100; i < 200; i++ {
		res, ok := c.Get(i)
		if !ok {
			t.Error(i, res, ok)
		}
	}
	t.Log(c.GetSet(1, 100))
	t.Log(c.Get(1))
	t.Log(c.GetSet(110, 110))
	t.Log(c.Get(110))
}
