package crawler

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"webcrawler/utils"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Urls map[url.URL]bool

type Crawler struct {
	Cfg          *Config
	UrlRetriever *UrlRetriever
	Logger       *log.Logger
}

type Config struct {
	MaxConnections                  int
	TimeOutConnection               int
	DelayAfterMaxConnectionsReached int
	RequestTimeout                  int
	Uris                            []string
	ShowResults                     bool
	LogLevel                        int
	Src                             string
}

type UrlRetriever struct {
	Urls Urls
}

type CrawlerInterface interface {
	GetUrls() (*Urls, error)
	Init()
	Crawl()
}

func (crawl Crawler) SetUrls(urls *Urls) error {
	urls, err := crawl.GetUrls()
	crawl.UrlRetriever.Urls = *urls
	return err
}

func (crawl Crawler) GetUrls() (*Urls, error) {
	urls := make(Urls)
	crawl.Logger.Debug("Trying to open file: %s", crawl.Cfg.Src)
	f, err := os.Open(crawl.Cfg.Src)
	if err != nil {
		return nil, err
	}
	crawl.Logger.Debug("Successfully opened file: %s", crawl.Cfg.Src)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, uri := range crawl.Cfg.Uris {
			url, err := utils.FormatUrl("https://"+scanner.Text(), uri)
			if err != nil {
				return nil, err
			}
			urls[url] = false
			url, err = utils.FormatUrl("https://www."+scanner.Text(), uri)
			if err != nil {
				return nil, err
			}
			urls[url] = false
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &urls, nil
}

func (crawl Crawler) cleanup() {
	if recovered := recover(); recovered != nil {
		crawl.Logger.Error("Failed to fetch url: ", recovered)
	}
}

func (crawl Crawler) FetchUrl(url url.URL, wg *sync.WaitGroup) (bool, error) {
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
		crawl.Logger.Error(err)
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(crawl.Cfg.RequestTimeout)*time.Second)
	defer cancel()
	crawl.Logger.Trace("Starting request to: %s", url.String())
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if err != nil {
		crawl.Logger.Error(err)
		panic(err)
	}
	defer resp.Body.Close()
	// Read data in bytes from the response
	_, err = io.ReadAll(resp.Body)
	// Get response stats code, and convert body of response to a readable string
	// Don't forget to close the response's body
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		crawl.Logger.Info("[ Exposed ] %s", url.String())
		crawl.UrlRetriever.Urls[url] = true
		return true, nil
	}
	return false, err
}

func (crawl Crawler) FetchUrlsConcurrently(urls Urls) {
	var wg sync.WaitGroup
	for url, _ := range urls {
		wg.Add(1)
		go crawl.FetchUrl(url, &wg)
	}
	wg.Wait()
}

func (crawl Crawler) FetchUrls() {
	var batchUrls = make(Urls, crawl.Cfg.MaxConnections)
	for urlToFetch, _ := range crawl.UrlRetriever.Urls {
		batchUrls[urlToFetch] = false
		if len(batchUrls) < crawl.Cfg.MaxConnections {
			continue
		}
		crawl.FetchUrlsConcurrently(batchUrls)
		crawl.Logger.Debugf("Sleeping for %d seconds", crawl.Cfg.DelayAfterMaxConnectionsReached)
		time.Sleep(time.Duration(crawl.Cfg.DelayAfterMaxConnectionsReached) * time.Second)
		crawl.Logger.Debugf("Finished sleeping after %d seconds", crawl.Cfg.DelayAfterMaxConnectionsReached)
		batchUrls = make(Urls, crawl.Cfg.MaxConnections)
	}
	// crawl the remanent urls
	if len(batchUrls) > 0 {
		crawl.Logger.Debugf("Crawling the rest of urls:", batchUrls)
		crawl.FetchUrlsConcurrently(batchUrls)
		batchUrls = make(Urls, crawl.Cfg.MaxConnections)
	}
}

func (crawl Crawler) GetVulnerableUrls() {
	for key, value := range crawl.UrlRetriever.Urls {
		if value {
			fmt.Printf("%s\n", key.String())
		}
	}
}

func (crawl Crawler) GetAllUrls() {
	for key, value := range crawl.UrlRetriever.Urls {
		fmt.Printf("%s:%t\n", key.String(), value)
	}
}

func (crawl Crawler) SetConfig() {
	flag.IntVar(&crawl.Cfg.LogLevel, "logLevel", 1, "Log level. Valid values from 1 to 6. Based on logrus levels.")
	flag.IntVar(&crawl.Cfg.MaxConnections, "concurrent", 150, "Max number of concurrent requests")
	flag.IntVar(&crawl.Cfg.RequestTimeout, "reqTimeout", 5, "Timeout (in seconds) before http request is aborted")
	flag.IntVar(&crawl.Cfg.TimeOutConnection, "connTimeout", 10, "Timeout (in seconds) before opening a new http connection")
	flag.IntVar(&crawl.Cfg.DelayAfterMaxConnectionsReached, "sleep", 0, "Timeout (in seconds) to sleep after the max number of concurrent connections has been reached")
	flag.BoolVar(&crawl.Cfg.ShowResults, "showResults", true, "Show all the sites that returned a valid response")
	flag.StringVar(&crawl.Cfg.Src, "file", "urls.txt", "Route of the file containing the urls to crawl, separated by newlines. Default to urls.txt")
	uris := []string{".env", ".git"}
	crawl.Cfg.Uris = uris
	flag.Parse()
}

func (crawl Crawler) SetLogger() {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.Level(crawl.Cfg.LogLevel))
	crawl.Logger = logger
}

func (crawl Crawler) SetInitialUrls() {
	urlMap := make(Urls)
	crawl.UrlRetriever.Urls = urlMap
}

func (crawl *Crawler) Init() {
	crawl.SetConfig()
	crawl.SetLogger()
	crawl.SetInitialUrls()
}

func (crawl Crawler) Crawl() {
	crawl.Logger.Debug("Starting the crawling process with the following configuration:", crawl.Cfg)
	crawl.Logger.Debugf("Getting the urls from: %s", crawl.Cfg.Src)
	urls, err := crawl.GetUrls()
	if err != nil {
		crawl.Logger.Fatalf("Failed to obtain the urls %s from: %s", crawl.Cfg.Src, err)
	}
	crawl.SetUrls(urls)
	crawl.FetchUrls()
	crawl.Logger.Debugf("Finished parsing the urls. %d urls to scan", len(crawl.UrlRetriever.Urls))
	crawl.Logger.Debug("Terminating the process.")
	crawl.Logger.Debug("Printing the results of the crawling process:")
	if crawl.Cfg.ShowResults {
		crawl.GetVulnerableUrls()
	}
}

func New() *Crawler {
	crawl := &Crawler{
		Cfg:          &Config{},
		UrlRetriever: &UrlRetriever{},
		Logger:       &log.Logger{},
	}
	crawl.Init()
	return crawl
}
