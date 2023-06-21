package firecache

import (
	"context"
	"sync"
	"time"
)

type cache struct {
	ctx      context.Context
	listener *listener
	cache    map[string]*dataCache
}

type dataCache struct {
	listener func(event ListenerEvent)
	timer    *time.Timer
	data     any
	datanil  bool
	err      error
}

func (c *cache) read(path string, query Q) (any, error) {
	key := parseKey(path, query)
	scope, ok := c.cache[key]
	if !ok {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		c.cache[key] = &dataCache{}
		scope = c.cache[key]
		scope.listener = func(event ListenerEvent) {
			if wg != nil {
				defer wg.Done()
			}

			if isDoc(path) {
				if event.Document == nil {
					scope.datanil = true
				}

				scope.data = event.Document
			} else {
				if event.DocumentList == nil {
					scope.datanil = true
				}

				scope.data = event.DocumentList
			}
		}
		c.listener.addListener(path, query, &scope.listener, func(err error) {
			if wg != nil {
				defer wg.Done()
			}
			scope.err = err
		})

		wg.Wait()
		wg = nil

		scope.timer = time.AfterFunc(time.Hour, func() {
			if scope == nil {
				return
			}

			c.listener.removeListener(path, query, &scope.listener)
			delete(c.cache, key)
		})
	} else {
		scope.timer.Reset(time.Hour)
	}

	if scope.err != nil {
		return nil, scope.err
	}

	if scope.datanil {
		return nil, &NoData{}
	}

	return scope.data, nil
}
