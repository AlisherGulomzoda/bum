package liberror

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ErrNonCustom = errors.New("abc")

func TestErrorEncoder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		err             error
		wantBody        string
		wantCode        int
		wantContentType string
	}{
		{
			name: "Error with json.Marshaller",
			err:  ErrBadRequest,
			wantBody: "{\"error\":\"bad request error\"," +
				"\"code\":\"BAD_REQUEST\"}\n",
			wantContentType: "application/json; charset=utf-8",
			wantCode:        http.StatusBadRequest,
		},
		{
			name: "Error without json.Marshaller",
			err:  ErrNonCustom,
			wantBody: "{\"error\":\"internal server error\"," +
				"\"code\":\"INTERNAL_ERROR\"}\n",
			wantContentType: "application/json; charset=utf-8",
			wantCode:        http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		p := tt
		t.Run(p.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx := context.Background()

			if err := ErrorEncoder(ctx, p.err, w); err != nil {
				t.Errorf("got error: %v", err)
			}

			if w.Code != p.wantCode {
				t.Errorf("Handlers.Get(), response code = %s(%d),"+
					" but want %s(%d)",
					http.StatusText(w.Code), w.Code,
					http.StatusText(p.wantCode), p.wantCode)
			}

			if got := w.Body.String(); got != p.wantBody {
				t.Errorf("Handlers.Get(), response body = %q, want %q", got, p.wantBody)
			}

			if got := w.Header().Get("Content-type"); got != p.wantContentType {
				t.Errorf("Handlers.Get(), response"+
					" content type = %q, want %q", got, p.wantContentType)
			}
		})
	}
}
