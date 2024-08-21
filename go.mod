module github.com/xdnv/ya-go-quiz

go 1.21.5

require go.uber.org/zap v1.27.0 // indirect

require (
	github.com/go-chi/chi/v5 v5.1.0
	internal/adapters/logger v1.0.0
	internal/app v1.0.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace internal/adapters/logger => ./internal/adapters/logger

replace internal/app => ./internal/app
