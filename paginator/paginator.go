package paginator

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/manat/link-header/link"
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
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   r.URL.Path,
	}

	totalPage := (total-1)/pageSize + 1 // Ceiling trick from https://stackoverflow.com/a/54006084

	return Paginator{
		PageSize:  pageSize,
		TotalPage: totalPage,
		URL:       url,
	}
}

// AsLinkHeader generates link header (RFC 5988) from the given page.
func (p *Paginator) AsLinkHeader(page int64) string {
	links := make([]link.Link, 0, 4)

	if page < p.TotalPage {
		p.setURLWithPageQuery(page + 1)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.NextRel,
		})

		p.setURLWithPageQuery(p.TotalPage)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.LastRel,
		})
	}

	if page > 1 {
		p.setURLWithPageQuery(page - 1)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.PrevRel,
		})

		p.setURLWithPageQuery(1)
		links = append(links, link.Link{
			URI: p.URL.String(),
			Rel: link.FirstRel,
		})
	}

	linkHeaders, err := link.Serialize(links)
	if err != nil {
		log.Fatal(err)
	}

	return linkHeaders
}

func (p *Paginator) setURLWithPageQuery(page int64) {
	pageQuery := p.URL.Query()
	pageQuery.Set(PageParam, strconv.Itoa(int(page)))

	p.URL.RawQuery = pageQuery.Encode()
}
