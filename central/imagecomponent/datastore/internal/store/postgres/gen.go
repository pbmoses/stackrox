package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ImageComponent --search-category IMAGE_COMPONENTS --postgres-migration-seq 4 --migrate-from "dackbox:image_component"
