package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"webcrawler/crawler"
)

func main() {
	var cfg crawler.Config
	var urls crawler.UrlParser
	uris := []string{".env", ".git"}
	cfg.Uris = uris
	flag.IntVar(&cfg.MaxConnections, "concurrent", 10, "Max number of concurrent requests")
	flag.IntVar(&cfg.RequestTimeout, "reqTimeout", 5, "Timeout (in seconds) before http request is aborted")
	flag.IntVar(&cfg.TimeOutConnection, "connTimeout", 10, "Timeout (in seconds) before opening a new http connection")
	flag.IntVar(&cfg.DelayAfterMaxConnectionsReached, "sleep", 0, "Timeout (in seconds) to sleep after the max number of concurrent connections has been reached")
	flag.StringVar(&urls.FileSrc, "file", "urls.txt", "Route of the file containing the urls to crawl, separated by newlines. Default to urls.txt")
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	flag.Parse()
	// Go package has many useful utilities for handling urls
	urlMap := make(map[url.URL]bool)
	urls.Urls = urlMap
	crwl := crawler.Crawler{
		Cfg:    &cfg,
		Logger: logger,
		Urls:   &urls,
	}
	crwl.Crawl()
}
