package firecache

func PaseDoc(data any) *Document {
	if data == nil {
		return nil
	}

	return data.(*Document)
}

func ParseColl(data any) *DocumentList {
	if data == nil {
		return nil
	}

	return data.(*DocumentList)
}
