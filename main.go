package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const url = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"

var priceRegexp = regexp.MustCompile("[^0-9.]+")

// NewProduct constructs a new Product from a byte array.
func NewProduct(resp []byte) (Product, error) {
	var p Product

	p.SetKBSize(len(resp))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return p, err
	}

	p.Title = doc.Find("h1").First().Text()
	ups := doc.Find(".pricePerUnit").First().Text()
	if err := p.ParseUnitPrice(ups); err != nil {
		return p, err
	}

	p.Description = doc.Find(".productDataItemHeader").First().Next().Text()

	return p, nil
}

// Product represents a parsed product HTML page.
type Product struct {
	Title       string  `json:"title"`
	Size        string  `json:"size"`
	UnitPrice   float32 `json:"unit_price"`
	Description string  `json:"description"`
}

// ParseUnitPrice parses and sets the unit price for a Product
// taking a string like "&pound;1.50" as an argument.
func (p *Product) ParseUnitPrice(price string) error {
	filtered := priceRegexp.ReplaceAllString(price, "")
	up, err := strconv.ParseFloat(filtered, 32)
	if err == nil {
		p.UnitPrice = float32(up)
	}

	return err
}

// SetKBSize sets the size of the response in KBs (eg, "53.2kb")
func (p *Product) SetKBSize(s int) {
	size := float64(s) / float64(1024)
	p.Size = strconv.FormatFloat(size, 'f', 2, 64) + "kb"
}

// Results represent a list of Products with an aggregated
// Total value of all products.
type Results struct {
	Products []Product `json:"results"`
	Total    float32   `json:"total"`
}

// AddProduct adds a new product to the result list and updates the total price.
func (r *Results) AddProduct(p Product) {
	// Add a product
	r.Products = append(r.Products, p)

	// TODO: consider recalculating total by
	// ranging over products
	r.Total += p.UnitPrice
}

// ProductURLs scrapes and returns a slice of links.
func ProductURLs(doc *goquery.Document) []string {
	var productURLs []string
	doc.Find(".productInfo h3 a").Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("href"); ok {
			productURLs = append(productURLs, url)
		}
	})

	return productURLs
}

func main() {

	indexDoc, err := goquery.NewDocument(url)
	fatal(err)

	results := new(Results)

	for _, url := range ProductURLs(indexDoc) {
		resp, err := http.Get(url)
		fatal(err)
		defer resp.Body.Close()

		respByte, err := ioutil.ReadAll(resp.Body)
		fatal(err)

		p, err := NewProduct(respByte)
		fatal(err)

		results.AddProduct(p)
	}

	data, err := json.MarshalIndent(results, "", "	")
	fmt.Print(string(data))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
