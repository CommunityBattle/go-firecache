package firecache

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var ctx context.Context = context.Background()

var firecacheOnce sync.Once
var firecacheInstance *Firecache

var listenerOnce sync.Once
var listenerInstance *listener

var cacheOnce sync.Once
var cacheInstance *cache

var databaseOnce sync.Once
var databaseInstance *database

func GetFirecache() *Firecache {
	if firecacheInstance == nil {
		firecacheOnce.Do(func() {
			firecacheInstance = &Firecache{}
			firecacheInstance.ctx = ctx
			firecacheInstance.listener = getListener()
			firecacheInstance.cache = getCache()
			firecacheInstance.database = getDatabase()
		})
	}

	return firecacheInstance
}

func getListener() *listener {
	if listenerInstance == nil {
		listenerOnce.Do(func() {
			listenerInstance = &listener{}
			listenerInstance.ctx = ctx
			listenerInstance.database = getDatabase()
			listenerInstance.cache = make(map[string]*listenerCache)
		})
	}

	return listenerInstance
}

func getCache() *cache {
	if cacheInstance == nil {
		cacheOnce.Do(func() {
			cacheInstance = &cache{}
			cacheInstance.ctx = ctx
			cacheInstance.listener = getListener()
			cacheInstance.cache = make(map[string]*dataCache)
		})
	}

	return cacheInstance
}

func getDatabase() *database {
	if databaseInstance == nil {
		databaseOnce.Do(func() {
			serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

			var client *firestore.Client
			var err error

			if len(serviceAccount) == 0 {
				client, err = firestore.NewClient(ctx, firestore.DetectProjectID)
			} else {
				sa := option.WithCredentialsFile(serviceAccount)
				client, err = firestore.NewClient(ctx, firestore.DetectProjectID, sa)
			}

			if err != nil {
				log.Fatalln(err)
			}

			databaseInstance = &database{}
			databaseInstance.ctx = ctx
			databaseInstance.firestore = client
		})
	}

	return databaseInstance
}
