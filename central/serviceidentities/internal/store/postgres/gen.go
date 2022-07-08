package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ServiceIdentity --get-all-func --postgres-migration-seq 54 --migrate-from "boltdb:service_identities"
