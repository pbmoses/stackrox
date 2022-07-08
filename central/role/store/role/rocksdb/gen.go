package rocksdb

//go:generate rocksdb-bindings-wrapper --type=Role --bucket=roles --cache --key-func GetName() --migrate-seq 50 --migrate-to roles
