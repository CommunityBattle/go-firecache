package firecache

import "cloud.google.com/go/firestore"

func ParseDocListener(data any) *firestore.DocumentSnapshot {
	if data == nil {
		return nil
	}

	return data.(*firestore.DocumentSnapshot)
}

func ParseCollListener(data any) *firestore.QuerySnapshot {
	if data == nil {
		return nil
	}

	return data.(*firestore.QuerySnapshot)
}

func PaseDoc(data any) Any {
	if data == nil {
		return nil
	}

	return data.(Any)
}

func ParseColl(data any) A {
	if data == nil {
		return nil
	}

	return data.(A)
}
