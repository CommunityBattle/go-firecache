package firecache

import (
	"context"
)

type cache struct {
	ctx      context.Context
	database *database
	listener *listener
}

func (c *cache) insert(path string, data any) (string, error) {
	return c.database.insert(path, data)
}

func (c *cache) update(path string, data U) error {
	return c.database.update(path, data)
}

func (c *cache) read(path string, query Q) (any, error) {
	return c.database.read(path, query)
}
