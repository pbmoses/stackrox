package postgres

//go:generate pg-table-bindings-wrapper --type=storage.Notifier --get-all-func --postgres-migration-seq 40 --migrate-from "boltdb:notifiers"
