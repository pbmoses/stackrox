package rocksdb

//go:generate rocksdb-bindings-wrapper --type=Pod --bucket=pods --track-index --migrate-seq 42 --migrate-to pods
