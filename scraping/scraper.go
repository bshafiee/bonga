package scraping

//Result represent a scraping listing
type Result struct {
	Title        string
	Price        string
	URL          string
	Img          string
	ID           string
	GeoTag       string
	Neighborhood string
}

type Parameter interface {
	String() string
}

type Scraper interface {
	Query([]Parameter) ([]Result, error)
}
