package firecache

type Query struct {
	field    string
	operator string
	value    string
}

type Q []Query
