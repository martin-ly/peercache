package peercache

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	writePath = "/peercache/write"
)

// Register registers a route that will accept (write) data into the cache.
// Peers' Write expects this route to exist on all remote Peers.
func Register(mux *http.ServeMux, cache *Cache) {
	mux.HandleFunc(writePath, func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		if key == "" {
			http.Error(w, "no key", http.StatusBadRequest)
			return
		}

		ttlStr := r.FormValue("ttl")
		if ttlStr == "" {
			ttlStr = "30s"
		}
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			http.Error(w, "bad ttl", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cache.Write(key, body, ttl)
	})
}

// Peers represents all peers in the cache group.
type Peers struct {
	urls []string
}

// NewPeers constructs a new Peers structure.
func NewPeers(urls []string) (*Peers, error) {
	validated := []string{}
	for _, u := range urls {
		v, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		validated = append(validated, v.String())
	}
	return &Peers{
		urls: validated,
	}, nil
}

// Write writes the
func (p *Peers) Write(key string, val []byte, ttl time.Duration) {
	for _, u := range p.urls {
		go func(url string) {
			_, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(val))
			if err != nil {
				log.Printf("peercache: %s: %s", url, err)
			}
		}(u + writePath + "?key=" + key + "&ttl=" + ttl.String())
	}
}
