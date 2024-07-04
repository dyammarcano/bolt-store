//go:build windows

package boltcache

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"testing"
)

// 100.000 registros escritos em 205,65 segundos = 486,2 registros por segundo na escrita

func TestNewRetrieveControlWrite(t *testing.T) {
	if err := NewBoltStore(context.TODO(), "testing"); err != nil {
		t.Errorf("error creating new BoltStore: %s", err)
	}

	if err := RegisterBucket("testing"); err != nil {
		t.Errorf("error registering bucket: %s", err)

	}

	count := 1000

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

func TestNewRetrieveControlRead(t *testing.T) {
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
