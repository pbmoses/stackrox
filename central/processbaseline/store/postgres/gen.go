package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ProcessBaseline --search-category PROCESS_BASELINES --postgres-migration-seq 45 --migrate-from "rocksdb:processWhitelists2"
