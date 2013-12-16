# rockgo

rockgo is a wrapper for [RocksDB](http://rocksdb.org)

## Building
You will need the shared library build of [RocksDB](https://github.com/facebook/rocksdb) installed on your machine. The current RocksDB version does not build the shared library by default, to build the shared library:
  
  make librocksdb.so

If the shared library and header files are not in a standard location for your platform, you will need to specify these:
  
  CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb/lib" go get github.com/wlattner/rockgo

## Notes
This package is incomplete. Additionally, the C api (used by rockgo) for RocksDB does not implement the full range of features available in the library.
