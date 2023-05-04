package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/fielderr"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	someData = map[int]any{
		1: "xd",
		2: 21,
		3: map[string]int{
			"xd": 1,
		},
	}
)

func TestCheckErr(t *testing.T) {
	tt := []struct {
		name   string
		err    error
		status int
		resp   interface{}
	}{
		{"unknown error", ErrUnknown, http.StatusBadRequest, nil},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeConflict), http.StatusConflict, nil},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeUnauthorized), http.StatusUnauthorized, someData},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srv := testServer(t, nil)
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			defer assert.NoError(t, r.Body.Close())
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			c := srv.engine.NewContext(r, w)
			if assert.NoError(t, srv.checkErr(c, tc.err)) {
				assert.Equal(t, tc.status, w.Code)
				wantBody, err := json.Marshal(tc.resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantBody), w.Body.String())
			}

		})
	}
}
