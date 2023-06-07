package firecache

type NoData struct{}

func (e *NoData) Error() string {
	return "no data"
}

type CollectionUsedForDocumentOperation struct{}

func (e *CollectionUsedForDocumentOperation) Error() string {
	return "collection used for document operation"
}

type AlreadyExists struct{}

func (e *AlreadyExists) Error() string {
	return "already exists"
}
