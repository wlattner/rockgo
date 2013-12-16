/*

Package rockgo provides a wrapper for the RocksDB C api. RocksDB is based on 
LevelDB.

  db := rockgo.New()
  db.Open("/path/to/db", true)  # true will create db if it does not already exist

The DB struct returned by New provides the standard Get, Put, Delete operations.
  
  err := db.Put([]byte("foo"), []byte("bar"))
  ...
  data, err := db.Get([]byte("foo"))
  ...
  err = db.Delete([]byte("foo"))
*/
package rockgo
