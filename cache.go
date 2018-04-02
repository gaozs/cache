package cache

import (
	"container/list"
	"errors"
	"sync"
)

// a cache can hold special num ID:Data pairs for quick use, like map, but will drop oldest ID:Data when full
// each R/W operation will set the ID:Data to newest
// ID is a type which can be used in map index, Data is interface{}

type CacheProvider struct {
	cacheLimitNum int
	count         int
	idxs          map[interface{}]*list.Element
	list          *list.List
	lock          sync.Mutex
}

type cacheData struct {
	id   interface{}
	data interface{}
}

func NewCache(cacheLimitNum int) (c *CacheProvider, err error) {
	if cacheLimitNum <= 0 {
		err = errors.New("Cache num is <=0!")
		return
	}
	cp := new(CacheProvider)
	cp.cacheLimitNum = cacheLimitNum
	cp.idxs = make(map[interface{}]*list.Element, cacheLimitNum)
	cp.list = list.New()
	c = cp
	return
}

func (c *CacheProvider) Set(id, data interface{}) (err error) {
	if c.count > c.cacheLimitNum {
		err = errors.New("Cache count is over limit, something is wrong!")
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	e, ok := c.idxs[id]
	if ok {
		// id in cache, update diretly
		e.Value.(*cacheData).data = data
		c.list.MoveToFront(e) // move this to front(which is newest r/w)
		return

	}
	err = c.addIDData(id, data)
	return
}

func (c *CacheProvider) Get(id interface{}) (data interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	e, ok := c.idxs[id]
	if ok {
		// id in cache
		data = e.Value.(*cacheData).data
		c.list.MoveToFront(e) // move this to front(which is newest r/w)
	}
	return
}

// get a data by id, if not exist set it
func (c *CacheProvider) GetSet(id, newData interface{}) (data interface{}, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	e, ok := c.idxs[id]
	if ok {
		// id in cache
		data = e.Value.(*cacheData).data
		c.list.MoveToFront(e) // move this to front(which is newest r/w)
		return
	}
	err = c.addIDData(id, newData)
	if err == nil {
		data = newData
	}
	return
}

func (c *CacheProvider) addIDData(id, data interface{}) (err error) {
	if c.count == c.cacheLimitNum {
		// cache is full, replace back element
		e := c.list.Back()
		if e == nil {
			err = errors.New("last element data is nil, something is wrong!")
			return
		}
		delete(c.idxs, e.Value.(*cacheData).id)
		c.idxs[id] = e
		e.Value.(*cacheData).id = id
		e.Value.(*cacheData).data = data
		c.list.MoveToFront(e)
		return
	}
	e := c.list.PushFront(&cacheData{id: id, data: data})
	c.idxs[id] = e
	c.count++
	return
}
