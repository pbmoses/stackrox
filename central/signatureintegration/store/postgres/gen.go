package postgres

//go:generate pg-table-bindings-wrapper --type=storage.SignatureIntegration --postgres-migration-seq 55 --migrate-from "rocksdb:signature_integrations"
