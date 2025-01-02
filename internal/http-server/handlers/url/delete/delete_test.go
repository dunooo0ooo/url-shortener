package delete_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		mockError  error
		statusCode int
		respError  string
	}{
		{
			name:       "Success",
			alias:      "test_alias",
			statusCode: http.StatusOK,
		},
		{
			name:       "Empty alias",
			alias:      "",
			statusCode: http.StatusNotFound,
			respError:  "invalid request",
		},
		{
			name:       "URL not found",
			alias:      "not_found_alias",
			mockError:  storage.ErrURLNotFound,
			statusCode: http.StatusOK,
		},
		{
			name:       "DeleteURL Error",
			alias:      "error_alias",
			mockError:  errors.New("unexpected error"),
			statusCode: http.StatusOK,
			respError:  "failed to delete url",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.alias != "" {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
			}

			handler := delete.Delete(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			r := chi.NewRouter()
			r.Delete("/delete/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/delete/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)

		})
	}
}
