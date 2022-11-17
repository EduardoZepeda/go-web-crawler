package crawler

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Crawler struct {
	Cfg    *Config
	Urls   *UrlParser
	Logger *log.Logger
}

type Config struct {
	MaxConnections                  int
	TimeOutConnection               int
	DelayAfterMaxConnectionsReached int
	RequestTimeout                  int
	Uris                            []string
}

type UrlParser struct {
	FileSrc string
	Urls    map[url.URL]bool
}

func (crawl *Crawler) AppendUrlToQueue(parsedUrl string, uri string) error {
	// According to go's documentation urls are recognized as [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	joinedUrl, err := url.JoinPath(parsedUrl, uri)
	if err != nil {
		return err
	}
	u, err := url.ParseRequestURI(joinedUrl)
	if err != nil {
		return err
	} else {
		crawl.Urls.Urls[*u] = false
	}
	return err
}

func (crawl *Crawler) GetUrls() error {
	f, err := os.Open(crawl.Urls.FileSrc)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, uri := range crawl.Cfg.Uris {
			err := crawl.AppendUrlToQueue("https://"+scanner.Text(), uri)
			if err != nil {
				return err
			}
			err = crawl.AppendUrlToQueue("https://www."+scanner.Text(), uri)
			if err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (crawl *Crawler) cleanup() {
	if recovered := recover(); recovered != nil {
		//crawl.Logger.Println("[ x ] Failed url: ", recovered)
	}
}

func (crawl *Crawler) ParseUrl(url url.URL, wg *sync.WaitGroup) (bool, error) {
	// Make sure to remove counter from waitgroup so crawler doesn't stop
	defer wg.Done()
	defer crawl.cleanup()
	// Create a new client to reuse connections,
	// Timeout default value: 10 seconds
	c := &http.Client{
		Timeout: time.Duration(crawl.Cfg.TimeOutConnection) * time.Second,
		// CheckRedirect prevents crawler to follow redirections, giving a false positive
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	// Create a New Request, it's not send at this point
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		crawl.Logger.Panic(err)
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(crawl.Cfg.RequestTimeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if err != nil {
		crawl.Logger.Panic(err)
		return false, err
	}
	defer resp.Body.Close()
	// Read data in bytes from the response
	_, err = io.ReadAll(resp.Body)
	// Get response stats code, and convert body of response to a readable string
	// Don't forget to close the response's body
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		crawl.Logger.Println("[ ok ] ", url.String())
		return true, nil
	}
	return false, err
}

func (crawl *Crawler) ParseUrlsConcurrently(urls *map[url.URL]bool) {
	var wg sync.WaitGroup
	for url, _ := range *urls {
		wg.Add(1)
		go crawl.ParseUrl(url, &wg)
	}
	wg.Wait()
}

func (crawl *Crawler) ParseUrls() {
	var batchUrls = make(map[url.URL]bool, crawl.Cfg.MaxConnections)
	for urlToFetch, _ := range crawl.Urls.Urls {
		batchUrls[urlToFetch] = false
		if len(batchUrls) < crawl.Cfg.MaxConnections {
			continue
		}
		crawl.ParseUrlsConcurrently(&batchUrls)
		time.Sleep(time.Duration(crawl.Cfg.DelayAfterMaxConnectionsReached) * time.Second)

		batchUrls = make(map[url.URL]bool, crawl.Cfg.MaxConnections)
	}
	// crawl the remanent urls
	if len(batchUrls) > 0 {
		crawl.ParseUrlsConcurrently(&batchUrls)
		batchUrls = make(map[url.URL]bool, crawl.Cfg.MaxConnections)
	}
}

func (crawler *Crawler) Crawl() {
	err := crawler.GetUrls()
	if err != nil {
		crawler.Logger.Fatalf("Failed to read the urls %s file: %s", crawler.Urls.FileSrc, err)
	}
	crawler.ParseUrls()
}
