package firecache

import (
	"fmt"
	"strings"
)

func parseKey(path string, query Q) string {
	key := path

	for _, cond := range query {
		key += fmt.Sprintf(":%s|%s|%v", cond.Field, cond.Operator, cond.Value)

		if cond.Order != nil {
			for _, order := range cond.Order {
				key += fmt.Sprintf("|%s|%v", order.By, order.Direction)
			}
		}

		if cond.Offset > 0 {
			key += fmt.Sprintf("|%v", cond.Offset)
		}

		if cond.Limit > 0 {
			key += fmt.Sprintf("|%v", cond.Limit)
		}
	}

	return key
}

func isDoc(path string) bool {
	hierarchy := strings.Split(path, "/")

	return len(hierarchy)%2 == 0
}

func parseDocIdFromPath(path string) string {
	hierarchy := strings.Split(path, "/")

	return hierarchy[len(hierarchy)-1]
}
