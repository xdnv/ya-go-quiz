package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"internal/adapters/logger"
	"internal/app"
	"internal/domain"
)

// universal storage
type UniStorage struct {
	config *app.ServerConfig
	ctx    context.Context
	//stor    *MemStorage
	db      *DbStorage
	timeout time.Duration
}

// init storage
func NewUniStorage(cf *app.ServerConfig) *UniStorage {

	//if cf.StorageMode == app.Database {
	var (
		conn *sql.DB
		err  error
	)

	conn, err = sql.Open("pgx", cf.DatabaseDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return &UniStorage{
		config:  cf,
		ctx:     context.Background(),
		db:      NewDbStorage(conn),
		timeout: 5 * time.Second,
	}
	// } else {
	// 	return &UniStorage{
	// 		config: cf,
	// 		ctx:    context.Background(),
	// 		stor:   NewMemStorage(),
	// 	}
	//}
}

func (t UniStorage) Bootstrap() error {
	//if t.config.StorageMode == app.Database {
	return t.db.Bootstrap(t.ctx)
	//}
	//return nil
}

func (t UniStorage) Close() {
	//if t.config.StorageMode == app.Database {
	t.db.Close()
	//}
}

func (t UniStorage) ToggleQuizAvailability(uuid string) error {
	return t.db.ToggleQuizAvailability(t.ctx, uuid)
}

func (t UniStorage) Ping() error {
	// if t.config.StorageMode == app.Database {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.Ping error"
	backoff := func(ctx context.Context) error {
		err := t.db.Ping(dbctx)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	//return t.db.Ping(dbctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}
	return err
	// } else {
	// 	return nil
	// }
}

func (t UniStorage) UpdateQuiz_(qd domain.QuizData) error {
	// if t.config.StorageMode == app.Database {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.UpdateQuiz error"
	backoff := func(ctx context.Context) error {
		err := t.db.UpdateQuiz_(dbctx, qd)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	//err := t.db.SetMetric(dbctx, name, metric)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}
	// } else {
	// 	t.stor.SetMetric(name, metric)
	// }
	return nil
}

func (t UniStorage) GetQuizRows(admin bool) (*[]domain.QuizRowData, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var qd *[]domain.QuizRowData

	errMsg := "UniStorage.GetQuizRows error"
	backoff := func(ctx context.Context) error {
		var err error
		qd, err = t.db.GetQuizRows(dbctx, admin)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return qd, err
}

func (t UniStorage) GetQuizData(uuid string) (*domain.QuizData, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var qd *domain.QuizData

	errMsg := "UniStorage.GetQuizData error"
	backoff := func(ctx context.Context) error {
		var err error
		qd, err = t.db.GetQuizData(dbctx, uuid)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return qd, err
}

// // Get a copy of Metric storage
// func (t UniStorage) GetMetrics() map[string]Metric {

// 	// Create the target map
// 	targetMap := make(map[string]Metric)

// 	// if t.config.StorageMode == app.Database {
// 	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
// 	defer cancel()

// 	errMsg := "UniStorage.GetMetrics error"
// 	backoff := func(ctx context.Context) error {
// 		var err error
// 		targetMap, err = t.db.GetMetrics(dbctx)
// 		return app.HandleRetriableDB(err, errMsg)
// 	}
// 	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
// 	//targetMap, err = t.db.GetMetrics(dbctx) //original w/o retries
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
// 		// return empty map
// 		return make(map[string]Metric)
// 	}
// 	return targetMap
// 	// } else {
// 	// 	// Get copy of original map
// 	// 	for key, value := range t.stor.Metrics {
// 	// 		targetMap[key] = value
// 	// 	}
// 	// }
// 	// return targetMap
// }

// func (t UniStorage) UpdateMetricS(mType string, mName string, mValue string) error {
// 	// if t.config.StorageMode == app.Database {
// 	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
// 	defer cancel()

// 	errMsg := "UniStorage.UpdateMetricS error"
// 	backoff := func(ctx context.Context) error {
// 		err := t.db.UpdateMetricS(dbctx, mType, mName, mValue)
// 		return app.HandleRetriableDB(err, errMsg)
// 	}
// 	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
// 	//return t.db.UpdateMetricS(dbctx, mType, mName, mValue)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
// 	}
// 	return err
// 	// } else {
// 	// 	return t.stor.UpdateMetricS(mType, mName, mValue)
// 	// }
// }

func (t UniStorage) UpdateQuiz(qd *domain.QuizData, errs *[]error) error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.UpdateQuiz error"
	backoff := func(ctx context.Context) error {
		err := t.db.UpdateQuiz(dbctx, qd, errs)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		*errs = append(*errs, err)
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}

	return nil
}
