package postgres

//go:generate pg-table-bindings-wrapper --type=Cluster --table=clusters  --uniq-key-func GetName()
