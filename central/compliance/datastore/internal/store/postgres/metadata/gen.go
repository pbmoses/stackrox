package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ComplianceRunMetadata --search-category COMPLIANCE_METADATA --migration-seq 18 --migrate-from rocksdb
