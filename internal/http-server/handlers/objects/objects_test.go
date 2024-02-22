package objects_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	resp "github.com/Zalozhnyy/sbt_k8s/internal/api/response"
	"github.com/Zalozhnyy/sbt_k8s/internal/http-server/handlers/objects"
	mapstorage "github.com/Zalozhnyy/sbt_k8s/internal/storage/map_storage"
	"github.com/Zalozhnyy/sbt_k8s/lib/slogdiscard"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestNewSaver(t *testing.T) {

	tests := []struct {
		name   string
		args   string
		id     string
		status int
		err    string
	}{
		{
			name:   "happy test",
			args:   `{"data": "{\"data\": \"json\"}"}`,
			status: http.StatusOK,
			err:    "",
			id:     "1",
		},
		{
			name:   "empty body",
			args:   ``,
			status: http.StatusBadRequest,
			err:    "empty request",
			id:     "1",
		},
		{
			name:   "not valid json",
			args:   `{"data": "{\"data\": \"json\""}`,
			status: http.StatusBadRequest,
			err:    "not valid input",
			id:     "1",
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			saver := mapstorage.New()

			r := chi.NewRouter()
			handler := objects.NewSaver(slogdiscard.NewDiscardLogger(), saver)
			r.Put("/objects/{id}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodPut, "/objects/"+tc.id, bytes.NewReader([]byte(tc.args)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()

			var resp resp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.err, resp.Error)

		})
	}
}

func TestNewGetter(t *testing.T) {
	tests := []struct {
		name   string
		getter func() objects.JsonGetter
		id     string
		err    string
		status int
	}{
		{
			name:   "happy test",
			getter: func() objects.JsonGetter { return mapstorage.New() },
			id:     "1",
			status: http.StatusBadRequest,
			err:    "id do not exists",
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			handler := objects.NewGetter(slogdiscard.NewDiscardLogger(), tc.getter())
			r.Get("/{id}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+tc.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()
			var resp resp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.err, resp.Error)

		})
	}
}
