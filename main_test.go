package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var productHTML = `
<html>
<body>
<div class="productSummary">
<div class="productTitleDescriptionContainer">
<h1>Sainsbury's Avocado Ripe & Ready XL Loose 300g</h1>
<p class="pricePerUnit">&pound;1.50<abbr title="per">/</abbr><abbr title="unit"><span class="pricePerUnitUnit">unit</span></abbr></p>
<h3 class="productDataItemHeader">Description</h3>
<div class="productText">
<p>Avocados</p>
<p></p>
</div>
</div>
</body>
</html>`

func TestProductJSON(t *testing.T) {
	data := []byte(productHTML)
	p, err := NewProduct(data)
	if err != nil {
		t.Error(err)
	}

	expSize := "0.40kb"
	if p.Size != expSize {
		t.Errorf("expected size %s, got %s", expSize, p.Size)
	}

	expTitle := "Sainsbury's Avocado Ripe & Ready XL Loose 300g"
	if p.Title != expTitle {
		t.Errorf("expected title \"%s\", got \"%s\"", expTitle, p.Title)
	}

	expPrice := float32(1.5)
	if p.UnitPrice != expPrice {
		t.Errorf("expected unit price %f, got %f", expPrice, p.UnitPrice)
	}

	expDesc := "\nAvocados\n\n"
	if p.Description != expDesc {
		t.Errorf("expected description \"%s\", got \"%s\"", expDesc, p.Description)
	}
}

var indexHTML = `
<html>
<body>
<div class="productInfo">
<h3><a href="foo-url">Avocados</a></h3>
<h3><a href="bar-url">Apples</a></h3>
<h3><a title="Kiwis">Kiwis</a></h3>
</div>
</body>
</html>`

func ExampleProductURLs() {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(indexHTML)))
	urls := ProductURLs(doc)
	fmt.Println(urls[0])
	fmt.Println(urls[1])
	fmt.Println(len(urls))
	// Output:
	// foo-url
	// bar-url
	// 2
}

func TestAddProduct(t *testing.T) {
	p1 := Product{UnitPrice: 1.5}
	p2 := Product{UnitPrice: 2.0}
	r := &Results{}
	r.AddProduct(p1)
	r.AddProduct(p2)

	if r.Total != 3.5 {
		t.Errorf("expected 3.5, got %f", r.Total)
	}
}
