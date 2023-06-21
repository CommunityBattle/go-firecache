package firecache

import (
	"context"

	"cloud.google.com/go/firestore"
)

type listener struct {
	ctx      context.Context
	database *database
	cache    map[string]*listenerCache
}

type listenerCache struct {
	callbacks    map[*func(event ListenerEvent)]func(event ListenerEvent)
	dataReceived bool
	initEvent    ListenerEvent
	unsubscribe  context.CancelFunc
}

func (l *listener) addListener(path string, query Q, callback *func(event ListenerEvent), errorhook func(error)) {
	key := parseKey(path, query)
	scope, ok := l.cache[key]
	if !ok {
		l.cache[key] = &listenerCache{}
		scope = l.cache[key]
		scope.callbacks = make(map[*func(event ListenerEvent)]func(event ListenerEvent))
		scope.dataReceived = false
		scope.unsubscribe = l.database.addListener(path, query, func(data any) {
			var event ListenerEvent

			if isDoc(path) {
				doc := data.(*firestore.DocumentSnapshot)
				docdata := Document(doc.Data())
				var docpointer *Document

				if doc.Exists() {
					docpointer = &docdata
				}

				scope.initEvent = ListenerEvent{Document: docpointer}
				event = ListenerEvent{Document: docpointer}
			} else {
				changes := data.(*firestore.QuerySnapshot).Changes
				docChangeList := make(DocumentChangeList, 0)
				for _, change := range changes {
					docChangeList = append(docChangeList, DocumentChangeEntry{Id: change.Doc.Ref.ID, Document: change.Doc.Data(), Kind: ChangeKind(change.Kind), NewIndex: change.NewIndex, OldIndex: change.OldIndex})
				}

				docs, err := data.(*firestore.QuerySnapshot).Documents.GetAll()
				if err != nil {
					errorhook(err)
					return
				}

				docList := make(DocumentList, 0)
				initDocChangeList := make(DocumentChangeList, 0)
				for index, doc := range docs {
					docList = append(docList, DocumentEntry{Id: doc.Ref.ID, Document: doc.Data()})
					initDocChangeList = append(initDocChangeList, DocumentChangeEntry{Id: doc.Ref.ID, Document: doc.Data(), Kind: Added, NewIndex: index, OldIndex: -1})
				}

				scope.initEvent = ListenerEvent{DocumentList: &docList, DocumentChangeList: &initDocChangeList}
				event = ListenerEvent{DocumentList: &docList, DocumentChangeList: &docChangeList}
			}

			for _, cb := range scope.callbacks {
				cb(event)
			}

			scope.dataReceived = true
		})
	}

	scope.callbacks[callback] = (*callback)

	if scope.dataReceived {
		(*callback)(scope.initEvent)
	}
}

func (l *listener) removeListener(path string, query Q, callback *func(event ListenerEvent)) {
	key := parseKey(path, query)
	scope, ok := l.cache[key]
	if !ok {
		return
	}

	delete(scope.callbacks, callback)

	if len(scope.callbacks) == 0 {
		scope.unsubscribe()
		delete(l.cache, key)
	}
}
