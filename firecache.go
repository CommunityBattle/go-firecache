package firecache

import (
	"context"
	"fmt"
	"runtime"
)

type Firecache struct {
	ctx      context.Context
	database *database
	cache    *cache
	listener *listener
}

func (f *Firecache) AddListener(path string, query Q, callback *func(event ListenerEvent), errorhook func(err error)) {
	f.listener.addListener(path, query, callback, errorhook)
}

func (f *Firecache) RemoveListener(path string, query Q, callback *func(event ListenerEvent)) {
	f.listener.removeListener(path, query, callback)
}

func (f *Firecache) Read(path string, query Q) (any, error) {
	return f.cache.read(path, query)
}

func (f *Firecache) ReadWithoutCache(path string, query Q) (any, error) {
	return f.database.read(path, query)
}

func (f *Firecache) Insert(path string, data any) (string, error) {
	return f.database.insert(path, data)
}

func (f *Firecache) Update(path string, data U) error {
	return f.database.update(path, data)
}

func (f *Firecache) Delete(path string, query Q) error {
	return f.database.delete(path, query)
}

func (f *Firecache) Monitor() {
	fmt.Println("--------------------------------------------------Firecache Monitor--------------------------------------------------")
	fmt.Println()
	fmt.Println("There are", runtime.NumGoroutine(), "goroutines running.")
	fmt.Println()
	fmt.Println("Currently active firecache caches:")
	for key, entry := range f.cache.cache {
		fmt.Println("\t", key+":", entry)
	}
	fmt.Println()
	fmt.Println("Currently active firecache listeners:")
	for key, entry := range f.listener.cache {
		fmt.Println("\t", key+":", entry)
	}
	fmt.Println("--------------------------------------------------Firecache Monitor--------------------------------------------------")
}
