package peercache_test

import (
	"testing"
	"time"

	"github.com/peterbourgon/peercache"
)

func TestTTL(t *testing.T) {
	c := peercache.NewCache(10 * time.Millisecond)
	c.Write("foo", []byte(`bar`), 25*time.Millisecond)

	data, ok := c.Read("foo")
	if !ok {
		t.Fatalf("expected data, but it wasn't there")
	}
	if expected, got := "bar", string(data); expected != got {
		t.Fatalf("expected '%s', got '%s'", expected, got)
	}

	time.Sleep(40 * time.Millisecond)

	if _, ok = c.Read("foo"); ok {
		t.Fatalf("after TTL, expected no data, but it was still there")
	}
}
