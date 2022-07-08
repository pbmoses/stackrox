package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ReportConfiguration --search-category REPORT_CONFIGURATIONS --postgres-migration-seq 47 --migrate-from "rocksdb:report_configs"
