package peercache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/peterbourgon/peercache"
)

func TestHandler(t *testing.T) {
	s1, c1 := newTestPeer(t)
	defer s1.Close()
	s2, c2 := newTestPeer(t)
	defer s2.Close()

	peers, err := peercache.NewPeers([]string{
		s1.URL,
		s2.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	peers.Write("foo", []byte(`bar`), 50*time.Millisecond)
	time.Sleep(25 * time.Millisecond)

	if val, ok := c1.Read("foo"); !ok {
		t.Fatalf("c1: expected successful read, but it was missing")
	} else if string(val) != "bar" {
		t.Fatalf("c1: expected 'bar', got '%s'", string(val))
	} else {
		t.Logf("c1: '%s' = %s OK", "foo", string(val))
	}

	if val, ok := c2.Read("foo"); !ok {
		t.Fatalf("c2: expected successful read, but it was missing")
	} else if string(val) != "bar" {
		t.Fatalf("c2: expected 'bar', got '%s'", string(val))
	} else {
		t.Logf("c2: '%s' = %s OK", "foo", string(val))
	}
}

func newTestPeer(t *testing.T) (*httptest.Server, *peercache.Cache) {
	h := http.NewServeMux()
	c := peercache.NewCache(10 * time.Millisecond)
	peercache.Register(h, c)
	s := httptest.NewServer(h)
	return s, c
}
