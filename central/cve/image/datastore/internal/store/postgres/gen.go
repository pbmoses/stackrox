package postgres

// //go:generate pg-table-bindings-wrapper --type=storage.CVE --table=image_cves --search-category VULNERABILITIES --postgres-migration-seq 18 --migrate-from "rocksdb:image_vuln"
