package firecache

import "context"

type cache struct {
	ctx      context.Context
	database *database
	listener *listener
}

func (c *cache) insert(path string, data any) (string, error) {
	return "", nil
}

func (c *cache) update(path string, data any) error {
	return nil
}

func (c *cache) read(path string, query Q) any {
	return c.database.read(path, query)
}
