package paginator

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/manat/go-link-header/link"
)

var (
	// PageParam stores default param value for page
	PageParam = "page"
	// PageSizeParam stores default param value for page size
	PageSizeParam = "page_size"
)

// Paginator represents functions necessary for generating pagination links.
type Paginator struct {
	PageSize  int64
	TotalPage int64
	URL       url.URL
}

// NewPaginator constructs a Paginator instance with default values.
func NewPaginator(r *http.Request, pageSize, total int64) Paginator {
	totalPage := (total-1)/pageSize + 1 // Ceiling trick from https://stackoverflow.com/a/54006084

	return Paginator{
		PageSize:  pageSize,
		TotalPage: totalPage,
		URL:       getURL(r),
	}
}

// AsLinkHeader generates link header (RFC 5988) from the given page.
// Empty string is returned when page is out of range.
func (p *Paginator) AsLinkHeader(page int64) string {
	if page < 1 || p.TotalPage < page {
		return ""
	}

	links := make([]link.Link, 0, 4)

	if page < p.TotalPage {
		nextPage := page + 1
		p.setURLWithPageQuery(page + 1)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.NextRel,
		})

		if nextPage != p.TotalPage {
			p.setURLWithPageQuery(p.TotalPage)
			links = append(links, link.Link{
				URI: p.URL.String(),
				Rel: link.LastRel,
			})
		}
	}

	if page > 1 {
		prevPage := page - 1
		p.setURLWithPageQuery(prevPage)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.PrevRel,
		})

		if prevPage > 1 {
			p.setURLWithPageQuery(1)
			links = append(links, link.Link{
				URI: p.URL.String(),
				Rel: link.FirstRel,
			})
		}
	}

	linkHeaders, err := link.Serialize(links)
	if err != nil {
		log.Println(err)
	}

	return linkHeaders
}

// setURLWithPageQuery mutates Paginator URL with page info.
func (p *Paginator) setURLWithPageQuery(page int64) {
	pageQuery := p.URL.Query()
	pageQuery.Set(PageParam, strconv.Itoa(int(page)))

	p.URL.RawQuery = pageQuery.Encode()
}

// getURL creates an URL from the given request.
func getURL(r *http.Request) url.URL {
	if r == nil {
		return url.URL{}
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	path, rawQuery := "", ""
	if r.URL != nil {
		path = r.URL.Path
		rawQuery = r.URL.RawQuery
	}

	return url.URL{
		Scheme:   scheme,
		Host:     r.Host,
		Path:     path,
		RawQuery: rawQuery,
	}
}
