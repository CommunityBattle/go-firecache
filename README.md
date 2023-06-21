# go-firecache
Go library for firestore caching module. 

This library is used to cache frequently requested data in a way, that stateless applications can perform the same request over and over without doing the actual request onto the database.
It uses snapshot listeners to stay in sync with the real firestore underneath.

Furthermore does it simplify the access to the firestore by wrapping the functionality into a CRUD pattern.

> **_NOTE:_** Altering the firestore in any way followed by an immediate read process wont show the altered state of the firestore. Use the `ReadWithoutCache` method to be able to immediate access the data.

## Usage
```go
package main

import (
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
	data, err := firecache.ReadWithoutCache("test_collection/new_document", nil)
	// data, err := firecache.Read("test_collection/new_document", nil)
	err = firecache.Update("test_collection/new_document", f.U{{Path: "foo", Value: "baz"}})
	err = firecache.Delete("test_collection/new_document", nil)

	firecache.RemoveListener("test_collection/new_document", nil, &handler)
	firecache.RemoveListener("test_collection", nil, &handler)
}
```

## Testing
This library can be tested through the test/e2e_test.go file.

The following steps are required to successfully run an E2E test:
- Create a new service account on the GCP platform that is permitted to CRUD the firestore
- Download the key set of the service account to the local machine
- Name the file "sa.json" and place it in the root folder of the library
- Run `make test`

There is also a possibility to monitor what the firecache is doing. This should not be used in production, however it also does not harm the production environment in any way.
 ```go
package main

import (
    f "github.com/CommunityBattle/go-firecache"
)

func main() {
    firecache := f.GetFirecache()
    firecache.Monitor()
}
 ```

 This will produce an output looking like this.

 ```zsh
--------------------------------------------------Firecache Monitor--------------------------------------------------

There are 9 goroutines running.

Currently active firecache caches:
         test/doc2: &{0x7a9e20 0xc00013c820 0xc00013a050 false}
         test: &{0x7a9e20 0xc0004963c0 0xc00049a360 false}

Currently active firecache listeners:
         test/doc2: &{map[0xc0004100f0:0x7a9e20] true {0xc00013a050 <nil> <nil>} 0x304a80}
         test: &{map[0xc0001096b0:0x7a9e20] true {<nil> 0xc00049a360 0xc00049a378} 0x304a80}
--------------------------------------------------Firecache Monitor--------------------------------------------------
 ```

## Limits
By design this library is not build to be lightning fast in any way. It is written efficiently but it will never be able to compete with nativly calling the firestore. 

Its cache ttl is hardcoded to 1 hour resetting it each time a request to the data is performed. This can lead to a lot of memory being used for the caching logic. It is recommended to use this only in manages environment. Furthermore it is recommended to not inflate the cache by caching each firestore call that is done. Consider what makes sense to be cached and only cache this data. Otherwise use the `ReadWithoutCache` method to interact with the firestore or use the go client directly.