package firecache

import "fmt"

func parseKey(path string, query Q) string {
	key := path

	for _, cond := range query {
		key += fmt.Sprintf(":%s|%s|%s", cond.field, cond.operator, cond.value)
	}

	return key
}
