package main

import (
	"fmt"

	f "github.com/CommunityBattle/go-firecache"
)

func main() {
	firecache := f.GetFirecache()

	handler := func(event f.ListenerEvent) {
		// do something
	}

	firecache.AddListener("test_collection/new_document", nil, &handler, nil)
	firecache.AddListener("test_collection", nil, &handler, func(err error) {
		// handle error
	})

	id, err := firecache.Insert("test_collection/new_document", f.Document{"foo": "bar"})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
	data, err := firecache.ReadWithoutCache("test_collection/new_document", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}
	// data, err := firecache.Read("test_collection/new_document", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(data)
	// }
	err = firecache.Update("test_collection/new_document", f.U{{Path: "foo", Value: "baz"}})
	if err != nil {
		fmt.Println(err)
	}
	err = firecache.Delete("test_collection/new_document", nil)
	if err != nil {
		fmt.Println(err)
	}

	firecache.RemoveListener("test_collection/new_document", nil, &handler)
}
