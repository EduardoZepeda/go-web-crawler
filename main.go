package main

import (
	"webcrawler/crawler"
)

func main() {
	var crwl crawler.CrawlerInterface
	crwl = crawler.New()
	crwl.Crawl()
}
