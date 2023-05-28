package firecache

import (
	"context"
)

type listener struct {
	ctx      context.Context
	database *database
	cache    map[string]*listenerCache
}

type listenerCache struct {
	callbacks    map[*func(data any)]func(data any)
	dataReceived bool
	data         any
	unsubscribe  context.CancelFunc
}

func (l *listener) addListener(path string, query Q, callback *func(data any)) error {
	key := parseKey(path, query)

	scope, ok := l.cache[key]

	if !ok {
		l.cache[key] = &listenerCache{}
		l.cache[key].callbacks = make(map[*func(data any)]func(data any))
		scope = l.cache[key]
	}

	scope.callbacks[callback] = (*callback)

	if scope.unsubscribe == nil {
		scope.unsubscribe = l.database.addListener(path, query, func(data any) {
			scope.data = data

			for _, cb := range scope.callbacks {
				cb(data)
			}
		})
	} else if scope.dataReceived {
		(*callback)(scope.data)
	}

	return nil
}

func (l *listener) removeListener(path string, query Q, callback *func(data any)) error {
	key := parseKey(path, query)

	scope, ok := l.cache[key]

	if !ok {
		return nil
	}

	delete(scope.callbacks, callback)

	if len(scope.callbacks) == 0 {
		scope.unsubscribe()
		delete(l.cache, key)
	}

	return nil
}
