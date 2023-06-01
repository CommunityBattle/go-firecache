# go-firecache
Go library for firestore caching module. 

This library is used to cache frequently requested data in a way, that stateless applications can perform the same request over and over without doing the actual request onto the database.
It uses snapshot listeners to stay in sync with the real firestore underneath.

Furthermore does it simplify the access to the firestore by wrapping the functionality into a CRUD pattern.

## Testing
This library can be tested through the test/e2e_test.go file.

The following steps are required to successfully run an E2E test:
- Create a new service account on the GCP platform that is permitted to CRUD the firestore
- Download the key set of the service account to the local machine
- Name the file "sa.json" and place it in the root folder of the library
- Run `make test`

## Usage
```go
package main

import (
    "github.com/CommunityBattle/go-firecache"
)

firecache := firecache.getInstance();

handler := func(data any) {
	//do something
}

err := firecache.AddListener("test_collection/new_document", nil, &handler)

_, err := firecache.Insert("test_collection/new_document", map[string]interface{}{"foo": "bar"});
data, err := firecache.Read("test_collection/new_document", nil)
err = firecache.Update("test_collection/new_document", map[string]interface{}{"foo": "baz"})
err := firecache.Delete("test_collection/new_document", nil)

firecache.RemoveListener("test_collection/new_document", nil, handler)
```