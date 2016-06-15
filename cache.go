// Package couchcache provides an implementation of httpcache.Cache that stores and
// retrieves data using Couchdb.
package couchcache

import (
	"log"
	"time"

	"github.com/cabify/go-couchdb"
)

// Cache objects store and retrieve data using Couchdb.
type Cache struct {
	// couch database where the cache will be stored
	Database *couchdb.DB
}

// New returns a new Cache
func New(db *couchdb.DB) *Cache {
	return &Cache{
		Database: db,
	}
}

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	result := record{}
	err := c.Database.Get(key, &result, nil)
	if err != nil {
		return []byte{}, false
	}
	return result.Content, true
}

func (c *Cache) Set(key string, content []byte) {
	rev, err := c.Database.Rev(key)
	if err != nil && !couchdb.NotFound(err) {
		return
	}
	_, err = c.Database.Put(key, &record{
		Created: time.Now(),
		Updated: time.Now(),
		Key:     key,
		Content: content,
	}, rev)
	if err != nil {
		log.Printf("Can't insert record in couch: %v\n", err)
		return
	}

	return
}

func (c *Cache) Delete(key string) {
	rev, err := c.Database.Rev(key)
	if err != nil && !couchdb.NotFound(err) {
		return
	}
	_, err = c.Database.Delete(key, rev)
	if err != nil {
		log.Printf("Can't remove record: %s", err)
	}
}

func (c *Cache) Indexes() {}

type record struct {
	Created time.Time
	Updated time.Time
	Key     string
	Content []byte
}
