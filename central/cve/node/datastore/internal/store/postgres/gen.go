package postgres

//go:generate pg-table-bindings-wrapper --type=storage.NodeCVE --table=node_cves --search-category NODE_VULNERABILITIES  --postgres-migration-seq 35 --migrate-from "rocksdb:image_vuln"
