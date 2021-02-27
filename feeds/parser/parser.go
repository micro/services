package parser

import (
	"errors"
	"net/http"
	"net/url"

	 "github.com/PuerkitoBio/goquery"
)

var (
	parsers = map[string]Parser{
		"cnbc.com": cnbcParser,
		"www.cnbc.com": cnbcParser,
	}
)

type Parser func(string) (string, error)

func Parse(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if v, ok := parsers[u.Host]; ok {
		return v(uri)
	}
	return "", errors.New("no parser for url")
}

func cnbcParser(url string) (string, error) {
        // Request the HTML page.
        res, err := http.Get(url)
        if err != nil {
		return "", err
        }

        defer res.Body.Close()

        if res.StatusCode != 200 {
                return "", errors.New("bad status code")
        }

        // Load the HTML document
        doc, err := goquery.NewDocumentFromReader(res.Body)
        if err != nil {
                return "", err
        }

        return doc.Find(".PageBuilder-col-9").Html()
}
