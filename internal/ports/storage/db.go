package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"strconv"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Main storage
type DbStorage struct {
	conn *sql.DB
}

// Init DB storage object
func NewDbStorage(conn *sql.DB) *DbStorage {
	return &DbStorage{conn: conn}
}

// Closes db connection
func (t DbStorage) Close() {
	t.conn.Close()
}

// prepare database
func (t DbStorage) Bootstrap(ctx context.Context) error {

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logger.Info("BOOTSTRAP STARTED")

	//check config
	//tableName := "public.config"
	dbKey := "DBVersion"
	dbVersion := "20240901"

	//Important! pgx does not support sql.Named(), use pgx.NamedArgs{} instead

	logger.Info("init config")

	//TODO: add version update procedure
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS public.config (
			key VARCHAR(128) NOT NULL PRIMARY KEY,
			value TEXT
		);
	`) //,
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO public.config (key, value)
			VALUES (@dbKey::text, @dbVersion::text)
		ON CONFLICT (key)
			DO UPDATE SET value = excluded.value;
	`,
		pgx.NamedArgs{
			"dbKey":     dbKey,
			"dbVersion": dbVersion,
		},
	)
	if err != nil {
		return err
	}

	logger.Info("init tests")

	// tests
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.tests (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
            ext_id VARCHAR(50) NOT NULL UNIQUE,
			version VARCHAR(8) NOT NULL,
			is_active BOOL NOT NULL,
			type VARCHAR(50) NOT NULL,
			name VARCHAR(250) NOT NULL,
			description VARCHAR(2048) NOT NULL			
        );
    `)
	if err != nil {
		return err
	}

	logger.Info("init questions")

	// questions
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.questions (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
			test_id UUID NOT NULL,
			ext_id VARCHAR(50) NOT NULL,
			type VARCHAR(50) NOT NULL,
			correct_value VARCHAR(150),
			text VARCHAR(1024) NOT NULL,
			comment VARCHAR(2048) NOT NULL,
			CONSTRAINT question_fk_test FOREIGN KEY (test_id) REFERENCES public.tests(id),
			CONSTRAINT question_uniq_test UNIQUE (ext_id, test_id)
        );
    `)
	if err != nil {
		return err
	}

	logger.Info("init options")

	// options
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.options (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
			question_id UUID NOT NULL,
			ext_id VARCHAR(50) NOT NULL,
			text VARCHAR(1024),
			value VARCHAR(50),
			is_correct BOOL,
			CONSTRAINT option_fk_question FOREIGN KEY (question_id) REFERENCES public.questions(id),
			CONSTRAINT option_uniq_question UNIQUE (ext_id, question_id),
			CONSTRAINT option_one_field_not_null CHECK (text IS NOT NULL OR value IS NOT NULL)
        );
    `)
	if err != nil {
		return err
	}

	logger.Info("init scores")

	// scores
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.scores (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
			test_id UUID NOT NULL,
			ext_id VARCHAR(50) NOT NULL,
			min_percent INT NOT NULL,
			max_percent INT NOT NULL,
			score INT NOT NULL,
			pass BOOL NOT NULL,
			comment VARCHAR(150),
			CONSTRAINT score_fk_test FOREIGN KEY (test_id) REFERENCES public.tests(id),
			CONSTRAINT score_uniq_test UNIQUE (ext_id, test_id)
        );
    `)
	if err != nil {
		return err
	}

	logger.Info("init results")

	// results
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.results (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
			test_id UUID NOT NULL,
			score_id UUID NOT NULL,
			pass_time TIMESTAMP NOT NULL, 
			result INT,
			score INT,
			is_passed BOOL,
			replies JSONB,
			CONSTRAINT result_fk_test FOREIGN KEY (test_id) REFERENCES public.tests(id),
			CONSTRAINT result_fk_score FOREIGN KEY (score_id) REFERENCES public.scores(id)
        );
    `)
	if err != nil {
		return err
	}

	logger.Info("BOOTSTRAP OK")

	// commit transaction
	return tx.Commit()
}

func (t DbStorage) Ping(ctx context.Context) error {
	return t.conn.PingContext(ctx)
}

func (t DbStorage) UpdateMetricS(ctx context.Context, mType string, mName string, mValue string) error {

	var val interface{}
	var err error
	query := ""

	switch mType {
	case "gauge":
		val, err = strconv.ParseFloat(mValue, 64)
		if err != nil {
			return err
		}
		query = `
		INSERT INTO public.gauges (id, value)
			VALUES (@id::text, @value::double precision)
		ON CONFLICT (id)
			DO UPDATE SET value = excluded.value;
	`
	case "counter":
		val, err = strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return err
		}
		query = `
		INSERT INTO public.counters (id, value)
			VALUES (@id::text, @value::bigint)
		ON CONFLICT (id)
			DO UPDATE SET value = public.counters.value + excluded.value;
	`
	default:
		return fmt.Errorf("unexpected metric type: %s", mType)
	}

	_, err = t.conn.ExecContext(ctx, query,
		pgx.NamedArgs{
			"id":    mName,
			"value": val,
		},
	)

	return err
}

func (t DbStorage) WriteQuizResult(ctx context.Context, qr domain.QuizResult) (string, error) {

	var result string

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO public.results (test_id, score_id, pass_time, result, score, is_passed, replies)
			VALUES (@test_id, @score_id, @pass_time, @result, @score, @is_passed, @replies)
		RETURNING id;
	`

	err = tx.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"test_id":   qr.TestID,
			"score_id":  qr.ScoreID,
			"pass_time": qr.PassTime,
			"result":    qr.Result,
			"score":     qr.Score,
			"is_passed": qr.IsPassed,
			"replies":   qr.Replies,
		},
	).Scan(&result)
	if err != nil {
		return result, err
	}

	qr.ID = result

	// commit transaction
	err = tx.Commit()
	return result, err
}

func (t DbStorage) UpdateQuizHeaderT(ctx context.Context, tx *sql.Tx, qd *domain.QuizData) (string, error) {

	var testID string

	query := `
		INSERT INTO public.tests (ext_id, version, is_active, type, name, description)
			VALUES (@ext_id, @version, @is_active, @type, @name, @description)
		ON CONFLICT (ext_id) DO UPDATE
			SET
				version = EXCLUDED.version,
				is_active = EXCLUDED.is_active,
				type = EXCLUDED.type,
				name = EXCLUDED.name,
				description = EXCLUDED.description
		RETURNING id;
	`
	err := tx.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"ext_id":      qd.ID,
			"version":     qd.Version,
			"is_active":   true,
			"type":        qd.Type,
			"name":        qd.Name,
			"description": qd.Description,
		},
	).Scan(&testID)

	return testID, err
}

func (t DbStorage) UpdateQuizQuestionsT(ctx context.Context, tx *sql.Tx, qd *domain.QuizData, errs *[]error) {

	query_question := `
		INSERT INTO public.questions (test_id, ext_id, type, text, comment)
			VALUES (@test_id, @ext_id, @type, @text, @comment)
		ON CONFLICT ON CONSTRAINT question_uniq_test DO UPDATE
			SET
				text = EXCLUDED.text,
				comment = EXCLUDED.comment
		RETURNING id;
	`

	query_option := `
		INSERT INTO public.options (question_id, ext_id, text, value, is_correct)
			VALUES (@question_id, @ext_id, @text, @value, @is_correct)
		ON CONFLICT ON CONSTRAINT option_uniq_question DO UPDATE
			SET
				text = EXCLUDED.text,
				value = EXCLUDED.value,
				is_correct = EXCLUDED.is_correct
		RETURNING id;
	`

	//process questions + options
	for _, q := range qd.Questions {
		var questionID string

		err := tx.QueryRowContext(ctx, query_question,
			pgx.NamedArgs{
				"test_id": qd.UUID,
				"ext_id":  q.ID,
				"type":    q.Type,
				"text":    q.Text,
				"comment": q.Comment,
			},
		).Scan(&questionID)

		if err != nil {
			*errs = append(*errs, err)
			continue
		}
		q.UUID = questionID

		for _, o := range q.Options {
			var optionID string

			err := tx.QueryRowContext(ctx, query_option,
				pgx.NamedArgs{
					"question_id": q.UUID,
					"ext_id":      o.ID,
					"text":        o.Text,
					"value":       o.Value,
					"is_correct":  o.IsCorrect,
				},
			).Scan(&optionID)

			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			o.UUID = optionID
		}
	}

}

func (t DbStorage) UpdateQuizScoresT(ctx context.Context, tx *sql.Tx, qd *domain.QuizData, errs *[]error) {

	query := `
		INSERT INTO public.scores (test_id, ext_id, min_percent, max_percent, score, pass, comment)
			VALUES (@test_id, @ext_id, @min_percent, @max_percent, @score, @pass, @comment)
		ON CONFLICT ON CONSTRAINT score_uniq_test DO UPDATE
			SET
				min_percent = EXCLUDED.min_percent,
				max_percent = EXCLUDED.max_percent,
				score = EXCLUDED.score,
				pass = EXCLUDED.pass,
				comment = EXCLUDED.comment
		RETURNING id;
	`

	//process scores
	for _, s := range qd.Scores {
		var scoreID string

		err := tx.QueryRowContext(ctx, query,
			pgx.NamedArgs{
				"test_id":     qd.UUID,
				"ext_id":      s.ID,
				"min_percent": s.MinPercent,
				"max_percent": s.MaxPercent,
				"score":       s.Score,
				"pass":        s.Pass,
				"comment":     s.Comment,
			},
		).Scan(&scoreID)

		if err != nil {
			*errs = append(*errs, err)
			continue
		}
		s.UUID = scoreID
	}
}

func (t DbStorage) UpdateQuiz(ctx context.Context, qd *domain.QuizData, errs *[]error) error {

	var testID string

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		*errs = append(*errs, err)
		return err
	}
	defer tx.Rollback()

	//process header
	testID, err = t.UpdateQuizHeaderT(ctx, tx, qd)
	if err != nil {
		*errs = append(*errs, err)
		// should not continue since header has insertion errors
		return err
	}

	qd.UUID = testID

	//process questions + options
	t.UpdateQuizQuestionsT(ctx, tx, qd, errs)

	//process scores
	t.UpdateQuizScoresT(ctx, tx, qd, errs)

	// commit transaction
	err = tx.Commit()
	return err
}

func (t DbStorage) GetQuizRows(ctx context.Context, admin bool) (*[]domain.QuizRowData, error) {

	var qd []domain.QuizRowData

	query := `
		SELECT id, ext_id, "version", is_active, "type", "name", description
		FROM public.tests
		WHERE is_active OR @adminMode;
	`

	rows, err := t.conn.QueryContext(ctx, query,
		pgx.NamedArgs{
			"adminMode": admin,
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("GetQuizRows: %s", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		qr := new(domain.QuizRowData)
		if err := rows.Scan(&qr.UUID, &qr.ID, &qr.Version, &qr.IsActive, &qr.Type, &qr.Name, &qr.Description); err != nil {
			logger.Error(fmt.Sprintf("GetQuizRows: %s", err))
			return nil, err
		}

		qr.WebID = domain.EncodeGUID(qr.UUID)
		qr.Link = "/quiz/" + domain.EncodeGUID(qr.UUID)

		qd = append(qd, *qr)
	}

	//logger.Info(fmt.Sprintf("GetQuizRows: got rows %s", len(qd))) //DEBUG

	return &qd, nil
}

func (t DbStorage) GetQuizData(ctx context.Context, uuid string) (*domain.QuizData, error) {

	var qd domain.QuizData
	var questionsJSON string

	query := `
	WITH question_data AS (
		select
			q.id AS id,
			q.ext_id AS ext_id,
			q.test_id AS test_id,
			q.text AS text,
			q."type" AS "type",
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'uuid', o.id,
					'id', o.ext_id,
					'text', o.text,
					'value', o.value,
					'is_correct', o.is_correct
				) ORDER BY o.ext_id ASC
			) AS options
		FROM public.questions q
		LEFT JOIN public.options o ON o.question_id = q.id
		GROUP BY q.id
	)
	select
		t.id,
		t.ext_id,
		t."version",
		t.is_active,
		t."type",
		t."name",
		t.description,
		COALESCE(
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'uuid', q.id,
					'id', q.ext_id,
					'type', q.type,
					'text', q.text,
					'options', q.options
				) ORDER BY q.ext_id ASC
			),
			'[]'
		) AS questions
	FROM public.tests t
	LEFT JOIN question_data q ON q.test_id = t.id
	WHERE t.id = @id
	GROUP BY t.id, t.ext_id, t."version", t.is_active, t."type", t."name", t.description;
	`

	err := t.conn.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"id": uuid,
		},
	).Scan(&qd.UUID, &qd.ID, &qd.Version, &qd.IsActive, &qd.Type, &qd.Name, &qd.Description, &questionsJSON)

	if err != nil {
		logger.Error(fmt.Sprintf("GetQuizData: %s", err))
		return nil, err
	}

	if err := json.Unmarshal([]byte(questionsJSON), &qd.Questions); err != nil {
		logger.Error(fmt.Sprintf("GetQuizData: %s", err))
		return nil, err
	}

	// for i := range qd.Questions {
	// 	logger.Info(fmt.Sprintf("GetQuizData: got option %v", qd.Questions[i].Options)) //DEBUG
	// }

	//logger.Info(fmt.Sprintf("GetQuizData: got quiz %s", qd)) //DEBUG

	return &qd, nil
}

func (t DbStorage) GetQuizScores(ctx context.Context, uuid string) (*[]domain.QuizScore, error) {

	var qs []domain.QuizScore

	query := `
		SELECT
			s.id AS id,
			s.ext_id AS ext_id,
			s.min_percent AS min_percent,
			s.max_percent AS max_percent,
			s.score AS score,
			s.pass AS pass,
			s.comment AS comment
		FROM public.scores s
		WHERE s.test_id = @id
		GROUP BY s.id
		ORDER BY s.ext_id;
	`

	rows, err := t.conn.QueryContext(ctx, query,
		pgx.NamedArgs{
			"id": uuid,
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("GetQuizScores: %s", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		qr := new(domain.QuizScore)
		if err := rows.Scan(&qr.UUID, &qr.ID, &qr.MinPercent, &qr.MaxPercent, &qr.Score, &qr.Pass, &qr.Comment); err != nil {
			logger.Error(fmt.Sprintf("GetQuizScores: %s", err))
			return nil, err
		}

		qs = append(qs, *qr)
	}

	return &qs, nil
}

func (t DbStorage) GetQuizResult(ctx context.Context, uuid string) (*domain.QuizResult, error) {

	var qr domain.QuizResult

	query := `
		SELECT
			r.id AS id,
			r.test_id AS test_id,
			r.score_id AS score_id,
			r.pass_time AS pass_time,
			r.result AS result,
			r.score AS score,
			r.is_passed AS is_passed
		FROM public.results r
		WHERE r.id = @id;
	`

	err := t.conn.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"id": uuid,
		},
	).Scan(&qr.ID, &qr.TestID, &qr.ScoreID, &qr.PassTime, &qr.Result, &qr.Score, &qr.IsPassed)

	if err != nil {
		logger.Error(fmt.Sprintf("GetQuizResult: %s", err))
		return nil, err
	}

	return &qr, nil
}

func (t DbStorage) ToggleQuizAvailability(ctx context.Context, uuid string) error {

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE public.tests
		SET is_active = NOT is_active
		WHERE id = @id;
	`

	_, err = tx.ExecContext(ctx, query,
		pgx.NamedArgs{
			"id": uuid,
		},
	)
	if err != nil {
		return err
	}

	// commit transaction
	err = tx.Commit()
	return err
}
