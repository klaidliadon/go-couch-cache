package couchcache

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/cabify/go-couchdb"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) Test(c *C) {
	trans := &http.Transport{
		MaxIdleConnsPerHost: 10,
		Proxy:               http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client, err := couchdb.NewClient("http://localhost:5984", trans)
	c.Assert(err, IsNil)
	cache := New(client.DB("cache"))
	cache.Indexes()

	key := "testKey"
	_, ok := cache.Get(key)

	c.Assert(ok, Equals, false)

	val := []byte("some bytes")
	cache.Set(key, val)

	retVal, ok := cache.Get(key)
	c.Assert(ok, Equals, true)
	c.Assert(string(retVal), Equals, string(val))

	val = []byte("some other bytes")
	cache.Set(key, val)

	retVal, ok = cache.Get(key)
	c.Assert(ok, Equals, true)
	c.Assert(string(retVal), Equals, string(val))

	cache.Delete(key)

	_, ok = cache.Get(key)
	c.Assert(ok, Equals, false)
}
