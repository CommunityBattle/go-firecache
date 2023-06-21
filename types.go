package firecache

import "cloud.google.com/go/firestore"

const (
	Asc  Direction = Direction(firestore.Asc)
	Desc Direction = Direction(firestore.Desc)

	Added    ChangeKind = ChangeKind(firestore.DocumentAdded)
	Modified ChangeKind = ChangeKind(firestore.DocumentModified)
	Removed  ChangeKind = ChangeKind(firestore.DocumentRemoved)
)

type Direction firestore.Direction
type ChangeKind firestore.DocumentChangeKind

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
	Offset   int
	Limit    int
}
type Q []Query

type Update firestore.Update
type U []Update

type Document map[string]interface{}
type DocumentEntry struct {
	Id       string
	Document Document
}
type DocumentList []DocumentEntry

type DocumentChangeEntry struct {
	Id       string
	Document Document
	Kind     ChangeKind
	OldIndex int
	NewIndex int
}
type DocumentChangeList []DocumentChangeEntry

type ListenerEvent struct {
	Document           *Document
	DocumentList       *DocumentList
	DocumentChangeList *DocumentChangeList
}
