//go:build windows

package boltcache

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"testing"
)

// 100.000 registros escritos em 205,65 segundos = 486,2 registros por segundo na escrita

func TestNewBoltStoreWrite(t *testing.T) {
	if err := NewBoltStore(context.TODO(), "testing"); err != nil {
		t.Errorf("error creating new BoltStore: %s", err)
	}

	if err := RegisterBucket("testing"); err != nil {
		t.Errorf("error registering bucket: %s", err)

	}

	count := 1_000_000

	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			SetString("testing", fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("testing %d", i)))))
		}()
	}

	wg.Wait()
}

// 100.000 registros lidos e descartados em 239,73 segundos = 417,1 registros por segundo na lida e descarte

func TestNewBoltStoreRead(t *testing.T) {
	count := Total("testing")

	for i := 0; i < count; i++ {
		v := Get("testing")
		if v == nil {
			t.Errorf("error getting value")
			return
		}

		v.Dispose()
	}

	if Total("testing") != 0 {
		t.Errorf("error on result")
	}
}

func TestNewBoltStoreBulkWrite(t *testing.T) {
	if err := NewBoltStore(context.TODO(), "testing"); err != nil {
		t.Errorf("error creating new BoltStore: %s", err)
	}
	defer func() {
		if err := Close(); err != nil {
			t.Errorf("error closing db: %s", err)
		}
		os.Remove(Path())
	}()

	if err := RegisterBucket("testing"); err != nil {
		t.Errorf("error registering bucket: %s", err)

	}

	count := 1_000_000
	items := make([]string, count)

	for i := 0; i < count; i++ {
		items[i] = fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("testing %d", i))))
	}

	if err := SetBulkString("testing", items); err != nil {
		t.Errorf("error setting bulk string: %s", err)
	}

	result := Total("testing")

	if result != count {
		t.Errorf("error on result: %d", result)
	}
}
