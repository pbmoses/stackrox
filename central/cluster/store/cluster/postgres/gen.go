package postgres

//go:generate pg-table-bindings-wrapper --type=storage.Cluster --search-category CLUSTERS --postgres-migration-seq 1 --migrate-from "rocksdb:clusters"
