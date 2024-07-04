package boltcache

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/segmentio/ksuid"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
)

var global *RetrieveControl

type KeyValue struct {
	db    *bolt.DB
	name  []byte
	Key   []byte
	Value []byte
}

func (kv *KeyValue) Dispose() error {
	return kv.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(kv.name)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		if err := bucket.Delete(kv.Key); err != nil {
			return fmt.Errorf("error deleting value: %s", err)
		}
		return nil
	})
}

type BucketName map[string][]byte

type RetrieveControl struct {
	db      *bolt.DB
	ctx     context.Context
	buckets BucketName
}

func NewRetrieveControl(ctx context.Context, name string) error {
	filePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.bolt", name))
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return err
	}

	global = &RetrieveControl{
		db:      db,
		ctx:     ctx,
		buckets: make(map[string][]byte),
	}
	return nil
}

func RegisterBucket(bucketName string) error {
	return global.registerBucket(bucketName)
}

func (r *RetrieveControl) registerBucket(bucketName string) error {
	b := sha1.Sum([]byte(bucketName))
	r.buckets[bucketName] = b[:]

	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(r.buckets[bucketName])
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return nil
	})
}

func Get(bucketName string) *KeyValue {
	return global.get(bucketName)
}

func (r *RetrieveControl) get(bucketName string) *KeyValue {
	var value *KeyValue
	r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(r.buckets[bucketName])
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		c := bucket.Cursor()
		k, v := c.First()
		value = &KeyValue{db: r.db, Key: k, Value: v, name: r.buckets[bucketName]}
		return nil
	})
	return value
}

func GetOnce(bucketName string) *KeyValue {
	return global.getOnce(bucketName)
}

func (r *RetrieveControl) getOnce(bucketName string) *KeyValue {
	var value *KeyValue
	r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(r.buckets[bucketName])
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		c := bucket.Cursor()
		k, v := c.First()
		value = &KeyValue{db: r.db, Key: k, Value: v, name: r.buckets[bucketName]}
		if err := bucket.Delete(k); err != nil {
			return fmt.Errorf("error deleting value: %s", err)
		}
		return nil
	})
	return value
}

func SetString(bucketName, value string) {
	global.set(bucketName, []byte(value))
}

func SetBytes(bucketName string, value []byte) error {
	return global.set(bucketName, value)
}

func (r *RetrieveControl) set(bucketName string, value []byte) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(r.buckets[bucketName])
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		if err := bucket.Put(ksuid.New().Bytes(), value); err != nil {
			return fmt.Errorf("error putting value: %s", err)
		}
		return nil
	})
}

func DeleteBucket(bucketName string) {
	global.deleteBucket(bucketName)
}

func (r *RetrieveControl) deleteBucket(bucketName string) error {
	if r.total(bucketName) > 0 {
		return fmt.Errorf("bucket not empty")
	}
	err := r.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(r.buckets[bucketName]); err != nil {
			return fmt.Errorf("error deleting bucket: %s", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting bucket: %s", err)
	}
	if err = r.db.Close(); err != nil {
		return fmt.Errorf("error closing db: %s", err)
	}

	return nil
}

func Total(bucketName string) int {
	return global.total(bucketName)
}

func (r *RetrieveControl) total(bucketName string) int {
	var total int
	r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(r.buckets[bucketName])
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			total++
		}
		return nil
	})
	return total
}

func GetAll(bucketName string) []string {
	return global.getAll(bucketName)
}

func (r *RetrieveControl) getAll(bucketName string) []string {
	var values = make([]string, 0)
	r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(r.buckets[bucketName])
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			values = append(values, string(v))
		}
		return nil
	})
	return values
}
