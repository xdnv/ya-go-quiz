module github.com/xdnv/ya-go-quiz

go 1.22.0

toolchain go1.22.7

require go.uber.org/zap v1.27.0 // indirect

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-chi/chi/v5 v5.1.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.23.0
	internal/adapters/logger v1.0.0
	internal/app v1.0.0
	internal/domain v1.0.0
	internal/ports/storage v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace internal/adapters/logger => ./internal/adapters/logger

replace internal/app => ./internal/app

replace internal/domain => ./internal/domain

replace internal/ports/storage => ./internal/ports/storage
