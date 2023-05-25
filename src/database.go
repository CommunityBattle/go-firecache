package firecache

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type database struct {
	ctx       context.Context
	firestore *firestore.Client
}

func (db *database) addListener(path string, query Q, callback func(data any)) func() {
	cancelCtx, unsubscribe := context.WithCancel(db.ctx)

	if isDoc(path) {
		it := db.firestore.Doc(path).Snapshots(cancelCtx)
		go listenDoc(it, callback)
	} else {
		it := db.firestore.Collection(path).Snapshots(cancelCtx)
		go listenColl(it, callback)
	}

	return unsubscribe
}

func (db *database) insert(path string, data any) (string, error) {
	var id string
	var err error

	if isDoc(path) {
		_, err = db.firestore.Doc(path).Set(db.ctx, data)
	} else {
		var doc *firestore.DocumentRef

		doc, _, err = db.firestore.Collection(path).Add(db.ctx, data)

		id = doc.ID
	}

	return id, err
}

func (db *database) update(path string, data any) error {
	return nil
}

func (db *database) read(path string, query Q) any {
	if isDoc(path) {
		doc, _ := db.firestore.Doc(path).Get(db.ctx)

		return doc.Data()
	} else {
		docs, _ := db.firestore.Collection(path).Documents(db.ctx).GetAll()

		var data []map[string]interface{}

		for _, doc := range docs {
			data = append(data, doc.Data())
		}

		return data
	}
}

func (db *database) delete(path string, query Q) error {
	return nil
}

func isDoc(path string) bool {
	hierarchy := strings.Split(path, "/")

	return len(hierarchy)%2 == 0
}

func listenDoc(iterator *firestore.DocumentSnapshotIterator, callback func(data any)) {
	for {
		snap, err := iterator.Next()

		if e := status.Code(err); e == codes.Canceled {
			return
		}

		callback(snap)
	}
}

func listenColl(iterator *firestore.QuerySnapshotIterator, callback func(data any)) {
	for {
		snap, err := iterator.Next()

		if e := status.Code(err); e == codes.Canceled {
			return
		}

		callback(snap)
	}
}
