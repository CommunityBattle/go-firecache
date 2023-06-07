package firecache

import "cloud.google.com/go/firestore"

type Direction int32

const (
	Asc  Direction = 1
	Desc Direction = 2
)

type Order struct {
	By        string
	Direction Direction
}

type O []Order

type Query struct {
	Field    string
	Operator string
	Value    interface{}
	Order    O
	Limit    int
}

type Q []Query

type Any map[string]interface{}

type A []Any

type Update firestore.Update

type U []Update
