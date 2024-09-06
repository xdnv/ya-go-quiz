package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// init
var _ = func() bool {

	testing.Init()
	return true
}()

// test GUID reversibility
func Test_EncodeGUID(t *testing.T) {
	type want struct {
		content string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "001 GUID reversibility test",
			want: want{
				content: "123e4567-e89b-12d3-a456-426655440000",
			},
			request: "123e4567-e89b-12d3-a456-426655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := tt.request

			tmp := EncodeGUID(request)

			result, err := DecodeGUID(tmp)
			require.NoError(t, err)

			assert.Equal(t, tt.want.content, result)
		})
	}

}
