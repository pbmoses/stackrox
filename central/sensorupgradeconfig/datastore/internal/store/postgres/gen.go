package postgres

//go:generate pg-table-bindings-wrapper --type=storage.SensorUpgradeConfig --singleton --postgres-migration-seq 52 --migrate-from "boltdb:sensor-upgrade-config"
