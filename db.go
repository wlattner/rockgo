package rockgo

// #cgo LDFLAGS: -lrocksdb
// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import "unsafe"
import "errors"

type DB struct {
	rocksDB          *C.rocksdb_t
	dbOpts           *C.rocksdb_options_t
	defaultReadOpts  *C.rocksdb_readoptions_t
	defaultWriteOpts *C.rocksdb_writeoptions_t
}

type ReadOptions *C.rocksdb_readoptions_t
type WriteOptions *C.rocksdb_writeoptions_t

// New returns a DB handle with default settings
// example usage:
//        db := rockgo.New()
//        db.Open("db", true)
func New() *DB {
	db := &DB{}
	db.dbOpts = C.rocksdb_options_create()
	db.defaultReadOpts = C.rocksdb_readoptions_create()
	db.defaultWriteOpts = C.rocksdb_writeoptions_create()
	return db
}

// Open opens a rocksdb database, if createIfMissing = true, Open will create
// the database if one is not found.
func (d *DB) Open(name string, createIfMissing bool) error {
	C.rocksdb_options_set_create_if_missing(d.dbOpts, boolToUchar(createIfMissing))

	var cErr *C.char
	defer C.free(unsafe.Pointer(cErr))

	dbCStr := C.CString(name)
	defer C.free(unsafe.Pointer(dbCStr))

	d.rocksDB = C.rocksdb_open(d.dbOpts, dbCStr, &cErr)
	if cErr != nil {
		err := C.GoString(cErr)
		return errors.New(err)
	}
	return nil
}

// Close closes a database, any calls made to DB after Close is called will
// cause problems, don't to it.
func (d *DB) Close() {
	C.rocksdb_close(d.rocksDB)
}

// Put sets key = value using default WriteOptions
func (d *DB) Put(key, value []byte) error {
	return d.PutWithOptions(key, value, d.defaultWriteOpts)
}

// PutWithOptions sets key = value using caller-supplied WriteOptions
func (d *DB) PutWithOptions(key, value []byte, opts WriteOptions) error {
	var kPtr, vPtr *C.char
	if len(key) != 0 {
		kPtr = (*C.char)(unsafe.Pointer(&key[0]))
	}
	if len(value) != 0 {
		vPtr = (*C.char)(unsafe.Pointer(&value[0]))
	}

	keyLen, valLen := len(key), len(value)
	var cErr *C.char
	defer C.free(unsafe.Pointer(cErr))
	C.rocksdb_put(d.rocksDB, opts, kPtr, C.size_t(keyLen), vPtr, C.size_t(valLen), &cErr)

	if cErr != nil {
		err := C.GoString(cErr)
		return errors.New(err)
	}
	return nil
}

// Get returns the value of key using default ReadOptions
func (d *DB) Get(key []byte) ([]byte, error) {
	return d.GetWithOptions(key, d.defaultReadOpts)
}

// GetWithOptions returns the value of key using caller-supplied ReadOptions
func (d *DB) GetWithOptions(key []byte, opts ReadOptions) ([]byte, error) {
	var kPtr *C.char
	if len(key) != 0 {
		kPtr = (*C.char)(unsafe.Pointer(&key[0]))
	}

	var cErr *C.char
	var cValLen C.size_t

	value := C.rocksdb_get(d.rocksDB, opts, kPtr, C.size_t(len(key)), &cValLen, &cErr)
	defer C.free(unsafe.Pointer(value))

	if cErr != nil {
		err := C.GoString(cErr)
		return nil, errors.New(err)
	}

	if value == nil {
		return nil, nil
	}
	return C.GoBytes(unsafe.Pointer(value), C.int(cValLen)), nil
}

// Delete removes a key using default WriteOptions
func (d *DB) Delete(key []byte) error {
	return d.DeleteWithOptions(key, d.defaultWriteOpts)
}

// DeleteWithOptions removes a key using caller-supplied WriteOptions
func (d *DB) DeleteWithOptions(key []byte, opts WriteOptions) error {
	var cErr, kPtr *C.char
	defer C.free(unsafe.Pointer(cErr))

	if len(key) != 0 {
		kPtr = (*C.char)(unsafe.Pointer(&key[0]))
	}

	C.rocksdb_delete(d.rocksDB, opts, kPtr, C.size_t(len(key)), &cErr)

	if cErr != nil {
		err := C.GoString(cErr)
		return errors.New(err)
	}
	return nil
}

func boolToUchar(b bool) C.uchar {
	uc := C.uchar(0)
	if b {
		uc = C.uchar(1)
	}
	return uc
}
