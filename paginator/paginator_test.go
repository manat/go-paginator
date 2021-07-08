package paginator

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsLinkHeader(t *testing.T) {
	t.Parallel()

	basicReq := &http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "goo.gl",
			Path:   "goo.gl/scx",
		},
	}

	fullReq := &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "goo.gl",
			Path:     "goo.gl/scx",
			RawQuery: "q2=test&page=1&page_size=10&q1=val",
		},
	}

	emptyReq := &http.Request{}

	tests := []struct {
		name  string
		pg    Paginator
		page  int64
		links string
	}{
		{
			"page 1/3 will show next and last",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(1), // page
			`<http://goo.gl/scx?page=2>; rel="next", <http://goo.gl/scx?page=3>; rel="last"`,
		},
		{
			"page 2/3 will show next and prev",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(2), // page
			`<http://goo.gl/scx?page=3>; rel="next", <http://goo.gl/scx?page=1>; rel="prev"`,
		},
		{
			"page 3/3 will show prev and first",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(3), // page
			`<http://goo.gl/scx?page=2>; rel="prev", <http://goo.gl/scx?page=1>; rel="first"`,
		},
		{
			"page 0/3 will show blank",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(0), // page
			``,
		},
		{
			"page -1/3 will show blank",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(-3), // page
			``,
		},
		{
			"page 4/3 will show blank",
			NewPaginator(
				basicReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(5), // page
			``,
		},
		{
			"page 1/3 with link header replaces original page and retains other query params",
			NewPaginator(
				fullReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(1), // page
			`<http://goo.gl/scx?page=2&page_size=10&q1=val&q2=test>; rel="next", <http://goo.gl/scx?page=3&page_size=10&q1=val&q2=test>; rel="last"`,
		},
		{
			"page 2/3 with link header replaces original page and retains other query params",
			NewPaginator(
				fullReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(2), // page
			`<http://goo.gl/scx?page=3&page_size=10&q1=val&q2=test>; rel="next", <http://goo.gl/scx?page=1&page_size=10&q1=val&q2=test>; rel="prev"`,
		},
		{
			"page 3/3 with link header replaces original page and retains other query params",
			NewPaginator(
				fullReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(3), // page
			`<http://goo.gl/scx?page=2&page_size=10&q1=val&q2=test>; rel="prev", <http://goo.gl/scx?page=1&page_size=10&q1=val&q2=test>; rel="first"`,
		},
		{
			"empty req will show blank",
			NewPaginator(
				emptyReq,
				int64(3),  // pageSize
				int64(7)), // total
			int64(5), // page
			``,
		},
		{
			"nil req will show blank",
			NewPaginator(
				nil,
				int64(3),  // pageSize
				int64(7)), // total
			int64(5), // page
			``,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.links, tt.pg.AsLinkHeader(tt.page))
		})
	}
}
