package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ServiceAccount --search-category SERVICE_ACCOUNTS --postgres-migration-seq 53 --migrate-from "rocksdb:service_accounts"
