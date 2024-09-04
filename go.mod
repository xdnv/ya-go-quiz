module github.com/xdnv/ya-go-quiz

go 1.21.5

require go.uber.org/zap v1.27.0 // indirect

require (
	github.com/go-chi/chi/v5 v5.1.0
	golang.org/x/crypto v0.17.0
	internal/adapters/logger v1.0.0
	internal/app v1.0.0
	internal/domain v1.0.0
	internal/ports/storage v1.0.0
)

require (
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace internal/adapters/logger => ./internal/adapters/logger

replace internal/app => ./internal/app

replace internal/domain => ./internal/domain

replace internal/ports/storage => ./internal/ports/storage
