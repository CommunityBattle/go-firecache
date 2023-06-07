package firecache

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type database struct {
	ctx       context.Context
	firestore *firestore.Client
}

func (db *database) addListener(path string, query Q, callback func(data any)) context.CancelFunc {
	cancelCtx, unsubscribe := context.WithCancel(db.ctx)

	if isDoc(path) {
		it := db.firestore.Doc(path).Snapshots(cancelCtx)
		go listenDoc(it, callback)
	} else {
		it := db.resolve(path, query).Snapshots(cancelCtx)
		go listenColl(it, callback)
	}

	return unsubscribe
}

func (db *database) insert(path string, data any) (string, error) {
	if isDoc(path) {
		doc, _ := db.firestore.Doc(path).Get(db.ctx)

		if doc.Exists() {
			return "", &AlreadyExists{}
		} else {
			_, err := db.firestore.Doc(path).Create(db.ctx, data)
			return "", err
		}
	} else {
		doc, _, err := db.firestore.Collection(path).Add(db.ctx, data)
		return doc.ID, err
	}
}

func (db *database) read(path string, query Q) (any, error) {
	if isDoc(path) {
		doc, _ := db.firestore.Doc(path).Get(db.ctx)

		if doc.Exists() {
			return doc.Data(), nil
		}

		return nil, &NoData{}
	} else {
		docs, err := db.resolve(path, query).Documents(db.ctx).GetAll()
		if err != nil {
			return nil, err
		}

		var data A

		for _, doc := range docs {
			data = append(data, doc.Data())
		}

		return data, nil
	}
}

func (db *database) update(path string, data U) error {
	if isDoc(path) {
		doc, _ := db.firestore.Doc(path).Get(db.ctx)

		if doc.Exists() {
			parsed := make([]firestore.Update, 0)

			for _, entry := range data {
				parsed = append(parsed, firestore.Update{
					Path:      entry.Path,
					FieldPath: entry.FieldPath,
					Value:     entry.Value,
				})
			}

			_, err := db.firestore.Doc(path).Update(db.ctx, parsed)
			return err
		} else {
			return &NoData{}
		}
	} else {
		return &CollectionUsedForDocumentOperation{}
	}
}

func (db *database) delete(path string, query Q) error {
	if isDoc(path) {
		doc, err := db.firestore.Doc(path).Get(db.ctx)
		if err != nil {
			return err
		}

		if doc.Exists() {
			collections, err := doc.Ref.Collections(db.ctx).GetAll()
			if err != nil {
				return err
			}

			for _, collection := range collections {
				docs, err := collection.Documents(db.ctx).GetAll()
				if err != nil {
					return err
				}

				for _, d := range docs {
					err := db.delete(path+"/"+d.Ref.ID, nil)
					if err != nil {
						return err
					}
				}
			}

			doc.Ref.Delete(db.ctx)
			return nil
		} else {
			return nil
		}
	} else {
		docs, err := db.resolve(path, query).Documents(db.ctx).GetAll()
		if err != nil {
			return err
		}

		for _, doc := range docs {
			err := db.delete(path+"/"+doc.Ref.ID, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (db *database) resolve(path string, query Q) firestore.Query {
	collectionRef := db.firestore.Collection(path)
	queryRef := collectionRef.Query

	if query != nil {
		for _, condition := range query {
			queryRef = queryRef.Where(condition.Field, condition.Operator, condition.Value)

			if condition.Order != nil {
				for _, order := range condition.Order {
					queryRef = queryRef.OrderBy(order.By, firestore.Direction(order.Direction))
				}
			}

			if condition.Limit > 0 {
				queryRef = queryRef.Limit(condition.Limit)
			}
		}
	}

	return queryRef
}

func isDoc(path string) bool {
	hierarchy := strings.Split(path, "/")

	return len(hierarchy)%2 == 0
}

func listenDoc(iterator *firestore.DocumentSnapshotIterator, callback func(data any)) {
	for {
		snap, err := iterator.Next()

		if e := status.Code(err); e == codes.Canceled {
			fmt.Println("canceled!!!")
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
