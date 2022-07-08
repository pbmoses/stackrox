package rocksdb

//go:generate rocksdb-bindings-wrapper --type=Secret --bucket=secrets --migrate-seq 51 --migrate-to secrets
