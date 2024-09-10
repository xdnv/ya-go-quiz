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

	var (
		conn *sql.DB
		err  error
	)

	if !cf.MockMode {
		conn, err = sql.Open("pgx", cf.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		conn = cf.MockConn
	}

	return &UniStorage{
		config:  cf,
		ctx:     context.Background(),
		db:      NewDbStorage(conn),
		timeout: 5 * time.Second,
	}
}

func (t UniStorage) Bootstrap() error {
	return t.db.Bootstrap(t.ctx)
}

func (t UniStorage) Close() {
	t.db.Close()
}

func (t UniStorage) ToggleQuizAvailability(uuid string) error {
	return t.db.ToggleQuizAvailability(t.ctx, uuid)
}

func (t UniStorage) Ping() error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.Ping error"
	backoff := func(ctx context.Context) error {
		err := t.db.Ping(dbctx)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}
	return err
}

func (t UniStorage) WriteQuizResult(qr domain.QuizResult) (string, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var result string

	errMsg := "UniStorage.WriteQuizResult error"
	backoff := func(ctx context.Context) error {
		var err error
		result, err = t.db.WriteQuizResult(dbctx, qr)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}

	return result, nil
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

func (t UniStorage) GetQuizScores(uuid string) (*[]domain.QuizScore, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var qs *[]domain.QuizScore

	errMsg := "UniStorage.GetQuizData error"
	backoff := func(ctx context.Context) error {
		var err error
		qs, err = t.db.GetQuizScores(dbctx, uuid)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return qs, err
}

func (t UniStorage) GetQuizResult(uuid string) (*domain.QuizResult, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var qr *domain.QuizResult

	errMsg := "UniStorage.GetQuizResult error"
	backoff := func(ctx context.Context) error {
		var err error
		qr, err = t.db.GetQuizResult(dbctx, uuid)
		return app.HandleRetriableDB(err, errMsg)
	}
	err := app.DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return qr, err
}

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
