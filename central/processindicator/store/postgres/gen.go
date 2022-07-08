package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ProcessIndicator --search-category PROCESS_INDICATORS --postgres-migration-seq 46 --migrate-from "rocksdb:process_indicators2"
