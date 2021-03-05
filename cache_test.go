package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c, err := NewCache(100)
	if err != nil {
		t.Fatal(err)
	}
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
	res, err := c.GetSet(1, 100)
	if err != nil {
		t.Error(err)
	}
	if res.(int) != 100 {
		t.Error(res, "!= 100")
	}
	err = c.Set(1, 200)
	if err != nil {
		t.Error(err)
	}
	var ok bool
	res, ok = c.Get(1)
	if !ok {
		t.Error("can't get 1")
	}
	if res.(int) != 200 {
		t.Error(res, "!= 200")
	}

	ok = c.Del(1)
	if !ok {
		t.Error("can't delete 1")
	}
	res, ok = c.Get(1)
	if ok {
		t.Error("should not get 1")
	}

	ok = c.Del(1)
	if ok {
		t.Error("can't repeat delete 1")
	}

	res, err = c.GetSet(150, 100)
	if err != nil {
		t.Error(err)
	}
	if res.(int) == 150 {
		t.Error(res, "== 150,should be 150*150")
	}
	t.Log(c.Usage())
}
