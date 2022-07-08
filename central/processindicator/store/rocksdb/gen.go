package rocksdb

//go:generate rocksdb-bindings-wrapper --type=ProcessIndicator --bucket=process_indicators2 --track-index --migrate-seq 46 --migrate-to process_indicators
