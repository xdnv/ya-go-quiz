package main

import (
	"database/sql/driver"
	"internal/app"
	"internal/ports/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
		},
		{
			name: "002 negative root test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				bodyHeader:  "",
			},
			request: "/bla",
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
	var testSc = app.ServerConfig{StorageMode: app.Database}
	testSc.MockMode = true
	testSc.Mock = &mock
	testSc.MockConn = db
	stor = storage.NewUniStorage(&testSc)

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

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
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
