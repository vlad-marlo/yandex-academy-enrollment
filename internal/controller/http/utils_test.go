package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
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
		{"unknown error", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
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
			if assert.NoError(t, srv.checkErr(c, "", tc.err)) {
				assert.Equal(t, tc.status, w.Code)
				wantBody, err := json.Marshal(tc.resp)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantBody), w.Body.String())
			}

		})
	}
}

func TestPaginationOpts_Limit(t *testing.T) {
	tt := []struct {
		name string
		opts *PaginationOpts
		want int
	}{
		{"nil", nil, 1},
		{"non nil", &PaginationOpts{limit: 10}, 10},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.opts.Limit())
		})
	}
}
func TestPaginationOpts_Offset(t *testing.T) {
	tt := []struct {
		name string
		opts *PaginationOpts
		want int
	}{
		{"nil", nil, 0},
		{"non nil", &PaginationOpts{offset: 10}, 10},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.opts.Offset())
		})
	}
}

func TestNewPaginationOpts(t *testing.T) {
	tt := []struct {
		name       string
		limit      string
		offset     string
		wantLimit  int
		wantOffset int
	}{
		{"main positive", "10", "12", 10, 12},
		{"not provided", "", "", 1, 0},
		{"provided only limit", "123", "", 123, 0},
		{"provided only offset", "", "123", 1, 123},
		{"limit not parsable", "bad", "123", 1, 123},
		{"offset not parsable", "123", "bad", 123, 0},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewPaginationOpts(tc.limit, tc.offset)
			if assert.NotNil(t, opts) {
				assert.Equal(t, tc.wantLimit, opts.limit)
				assert.Equal(t, tc.wantOffset, opts.offset)
			}
		})
	}
}
