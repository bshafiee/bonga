package craiglist

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bshafiee/bonga/scraping"
)

const baseURL = "https://vancouver.craigslist.org/search/sfc/apa?bundleDuplicates=1&s="

type craiglistScraper struct {
	client *http.Client
}

func NewCraiglistScraper() scraping.Scraper {
	return &craiglistScraper{}
}

func (*craiglistScraper) fetch(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"+strconv.Itoa(rand.Int()))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}

func (cs *craiglistScraper) Query(params []scraping.Parameter) ([]scraping.Result, error) {
	pageCounter := 0
	results := make([]scraping.Result, 0)
	unique := make(map[string]bool)

	for {
		//1) build the URL
		url := baseURL + strconv.Itoa(pageCounter)
		for i, p := range params {
			if i == 0 {
				url += "&"
			}
			url += p.String()
			if i != (len(params) - 1) {
				url += "&"
			}
		}
		//2) query it
		doc, err := cs.fetch(url)
		fmt.Println(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch results from craiglist:%s", err)
		}
		//3) parse itx
		resRows := doc.Find(".result-row")
		if resRows.Length() == 0 {
			break
		}
		resRows.Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			id := s.AttrOr("data-pid", "")
			title := s.Find(".result-title").First().Text()
			price := s.Find(".result-price").First().Text()
			url := s.Find("a").AttrOr("href", "")
			img := s.Find(".swipe-wrap div img").AttrOr("src", "")
			date := s.Find(".result-date").First().Text()
			if !unique[title] {
				unique[title] = true
				results = append(results, scraping.Result{
					URL:   url,
					Title: title,
					Price: price,
					Img:   img,
					ID:    id,
					Date:  date,
				})
			}
		})
		pageCounter += resRows.Length()
		time.Sleep(time.Second * 3)
	}
	return results, nil
}
