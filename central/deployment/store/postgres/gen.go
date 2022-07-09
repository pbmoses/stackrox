package postgres

//go:generate pg-table-bindings-wrapper --type=storage.Deployment --search-category DEPLOYMENTS --references=storage.Image,namespaces:storage.NamespaceMetadata --postgres-migration-seq 3 --migrate-from "dackbox:deployments"
