package rocksdb

//go:generate rocksdb-bindings-wrapper --type=ReportConfiguration --bucket=report_configs --cache --key-func GetId() --migrate-seq 47 --migrate-to report_configurations
