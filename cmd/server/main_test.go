package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"internal/app"
	"internal/ports/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// converts pgx.NamedArgs to sqlmock compatible type -- int
type PgxCustomConverter struct{}

func (c PgxCustomConverter) ConvertValue(v interface{}) (driver.Value, error) {
	if vt, ok := v.(pgx.NamedArgs); ok {
		//return count of params to check for completeness
		return len(vt), nil
	}
	return v, nil
}

// User-defined sqlmock Match
type AnyArg struct{}

func (a AnyArg) Match(v driver.Value) bool {
	return true
}

// common functions
func NewTestServerConfig(db *sql.DB, mock *sqlmock.Sqlmock) *app.ServerConfig {
	var testSc = app.ServerConfig{StorageMode: app.Database}
	testSc.MockMode = true
	testSc.Mock = mock
	testSc.MockConn = db

	return &testSc
}

// init
var _ = func() bool {
	testing.Init()
	return true
}()

func Test_index(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		bodyHeader  string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name: "001 positive root test",
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  200,
				bodyHeader:  "<html>",
			},
			request: "/",
			method:  http.MethodGet,
		},
		{
			name: "002 negative root test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				bodyHeader:  "",
			},
			request: "/bla",
			method:  http.MethodGet,
		},
	}

	// create db mock with custom pgxNamedArgs converter
	converter := PgxCustomConverter{}
	db, mock, err := sqlmock.New(sqlmock.ValueConverterOption(converter))

	if err != nil {
		t.Fatalf("error opening sqlmock: '%s'", err)
	}
	defer db.Close()

	//switch sever config to mock mode
	testSc := NewTestServerConfig(db, &mock)
	stor = storage.NewUniStorage(testSc)

	// Expect query
	mock.ExpectQuery(`
	SELECT id, ext_id, "version", is_active, "type", "name", description
	FROM public.tests
	WHERE is_active OR @adminMode;
`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ext_id", "version", "is_active", "type", "name", "description"}).
			AddRow(uuid.New(), "TST001", "20240901", true, "go-quiz-test", "Test DMO", "Test DESC"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			index(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			//bodyString := string(bodyBytes)
			//assert.True(t, strings.Contains(bodyString, tt.want.bodyHeader))
		})
	}
}

func Test_handleResults(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		bodyHeader  string
	}
	tests := []struct {
		name      string
		request   string
		parameter string
		method    string
		want      want
	}{
		{
			name: "001 positive result test",
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  200,
				bodyHeader:  "<html>",
			},
			request:   "/results/{id}",
			parameter: "qMp_rJ4dH97-mx9jdsmFkvP",
			method:    http.MethodGet,
		},
		{
			name: "002 negative result test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				bodyHeader:  "",
			},
			request:   "/results/{id}",
			parameter: "56789",
			method:    http.MethodGet,
		},
	}

	// create db mock with custom pgxNamedArgs converter
	converter := PgxCustomConverter{}
	db, mock, err := sqlmock.New(sqlmock.ValueConverterOption(converter))

	if err != nil {
		t.Fatalf("error opening sqlmock: '%s'", err)
	}
	defer db.Close()

	//switch sever config to mock mode
	testSc := NewTestServerConfig(db, &mock)
	stor = storage.NewUniStorage(testSc)

	//uuid := uuid.New()

	// Expect query
	mock.ExpectQuery(`
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
`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "test_id", "score_id", "pass_time", "result", "score", "is_passed"}).
			AddRow(
				"123e4567-e89b-12d3-a456-426655440000",
				"123e4567-e89b-12d3-a456-426655440000",
				"123e4567-e89b-12d3-a456-426655440000",
				time.Now(),
				75,
				3,
				true))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.request, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.parameter)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handleResults(w, r)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			//bodyString := string(bodyBytes)
			//assert.True(t, strings.Contains(bodyString, tt.want.bodyHeader))
		})
	}
}
