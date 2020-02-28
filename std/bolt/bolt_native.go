package bolt

import (
	"bytes"
	. "github.com/candid82/joker/core"
	bolt "go.etcd.io/bbolt"
	"os"
	"unsafe"
)

type (
	// TODO: wrapper types like this can probably be auto generated
	BoltDB struct {
		*bolt.DB
		hash uint32
	}
)

var boltDBType *Type

func MakeBoltDB(db *bolt.DB) BoltDB {
	res := BoltDB{db, 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(db)))
	return res
}

func (db BoltDB) ToString(escape bool) string {
	return "#object[BoltDB]"
}

func (db BoltDB) Equals(other interface{}) bool {
	if otherDb, ok := other.(BoltDB); ok {
		return db.DB == otherDb.DB
	}
	return false
}

func (db BoltDB) GetInfo() *ObjectInfo {
	return nil
}

func (db BoltDB) GetType() *Type {
	return boltDBType
}

func (db BoltDB) Hash() uint32 {
	return db.hash
}

func (db BoltDB) WithInfo(info *ObjectInfo) Object {
	return db
}

func EnsureBoltDB(args []Object, index int) BoltDB {
	switch c := args[index].(type) {
	case BoltDB:
		return c
	default:
		panic(RT.NewArgTypeError(index, c, "BoltDB"))
	}
}

func ExtractBoltDB(args []Object, index int) *bolt.DB {
	return EnsureBoltDB(args, index).DB
}

func open(filename string, mode int) *bolt.DB {
	RT.GIL.Unlock()
	db, err := bolt.Open(filename, os.FileMode(mode), nil)
	RT.GIL.Lock()
	PanicOnErr(err)
	return db
}

func close(db *bolt.DB) Nil {
	RT.GIL.Unlock()
	err := db.Close()
	RT.GIL.Lock()
	PanicOnErr(err)
	return NIL
}

func createBucket(db *bolt.DB, name string) Nil {
	db.Update(func(tx *bolt.Tx) error {
		RT.GIL.Unlock()
		_, err := tx.CreateBucket([]byte(name))
		RT.GIL.Lock()
		PanicOnErr(err)
		return nil
	})
	return NIL
}

func deleteBucket(db *bolt.DB, name string) Nil {
	db.Update(func(tx *bolt.Tx) error {
		RT.GIL.Unlock()
		err := tx.DeleteBucket([]byte(name))
		RT.GIL.Lock()
		PanicOnErr(err)
		return nil
	})
	return NIL
}

func getBucket(tx *bolt.Tx, bucket string) *bolt.Bucket {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		panic(RT.NewError("Bucket doesn't exists: " + bucket))
	}
	return b
}

func nextSequence(db *bolt.DB, bucket string) Int {
	var id uint64
	db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, bucket)
		RT.GIL.Unlock()
		id, _ = b.NextSequence()
		RT.GIL.Lock()
		return nil
	})
	return MakeInt(int(id))
}

func put(db *bolt.DB, bucket, key, value string) Nil {
	db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, bucket)
		RT.GIL.Unlock()
		err := b.Put([]byte(key), []byte(value))
		RT.GIL.Lock()
		PanicOnErr(err)
		return nil
	})
	return NIL
}

func delete(db *bolt.DB, bucket, key string) Nil {
	db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, bucket)
		RT.GIL.Unlock()
		err := b.Delete([]byte(key))
		RT.GIL.Lock()
		PanicOnErr(err)
		return nil
	})
	return NIL
}

func get(db *bolt.DB, bucket, key string) Object {
	var v []byte
	db.View(func(tx *bolt.Tx) error {
		b := getBucket(tx, bucket)
		RT.GIL.Unlock()
		v = b.Get([]byte(key))
		RT.GIL.Lock()
		return nil
	})
	if v == nil {
		return NIL
	}
	return MakeString(string(v))
}

func byPrefix(db *bolt.DB, bucket, prefix string) *Vector {
	res := EmptyVector()
	db.View(func(tx *bolt.Tx) error {
		c := getBucket(tx, bucket).Cursor()
		pr := []byte(prefix)
		RT.GIL.Unlock()
		for k, v := c.Seek(pr); k != nil && bytes.HasPrefix(k, pr); k, v = c.Next() {
			res = res.Conjoin(NewVectorFrom(MakeString(string(k)), MakeString(string(v))))
		}
		RT.GIL.Lock()
		return nil
	})
	return res
}

func init() {
	boltDBType = RegType("BoltDB", (*BoltDB)(nil), "Wraps Bolt DB type")
}

func initNative() {
}
