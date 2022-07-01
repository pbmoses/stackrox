package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ImageCVE --table=image_cves --search-category IMAGE_VULNERABILITIES  --postgres-migration-seq 18 --migrate-from "rocksdb:image_vuln"
