package firecache

type Order struct {
	by        string
	direction string
}

type O []Order

type Query struct {
	field    string
	operator string
	value    string
	order    O
	limit    int
}

type Q []Query
