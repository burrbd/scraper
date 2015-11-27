### Product scraper

Application follows and scrapes products from the following URL:
```http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html```

Application uses https://github.com/PuerkitoBio/goquery project for parsing HTML contents.

#### Installation

1. ```git clone https://github.com/burrbd/scraper.git```
1. ```cd``` into project dir
1. Run ```go get && go build```

#### Usage

Run ```./scraper``` to print results to stdout, or ```./scraper > output.json``` to print results to file.

Tests can be run with ```go test```.
