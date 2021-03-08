package paginator

import (
	"net/http"
	"testing"
)

func TestAsLinkHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		pg    Paginator
		links string
		err   error
	}{
		{
			"ok for ChannelItem object with all fields",
			NewPaginator(
				&http.Request{
					Host: "",
				},
				int64(3),
				int64(7)),
			"",
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

		})
	}
}
