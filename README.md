# Go crawler for .env and .git

This crawler was inspired by [scanning-26-million-domains-for-exposed-env-files](https://hackernoon.com/scanning-26-million-domains-for-exposed-env-files) article. It uses concurrency to crawl a list of domains and check for exposed .env and .git uris, in plain or www subdomains.

This project uses Go 1.19. Maybe I'll add a Docker file later.

## Quickstart

### You need a file with urls

You need a file with domain names separated by newlines. Default file name is *urls.txt* at the root of the project.

```bash
example.org
domain.org
 ```
 
 The crawler will request the following urls:
 
 ```bash
https://www.example.org/.env
https://www.example.org/.git
https://example.org/.env
https://example.org/.git
# ...
 ```

 ### Run the crawler

 Once you got the file set, you can start crawling with:

 ```bash
go run main.go
 ```
 
 ## Acceptance criteria
 
 The crawler consider a successful response as one with a code status between 200 and 300 and ignores any redirection. I'm aware this could lead to some false positives. Please modify the code according to your own needs.
 
 ## Flag options
 
 The following options can be used to customize the crawler behaviour.
 
 - concurrent: Max number of concurrent requests. Default to 10
 - reqTimeout: Timeout (in seconds) before http request is aborted. Default to 5
 - connTimeout: Timeout (in seconds) before opening a new http connection. Default to 10
 - sleep: Timeout (in seconds) to sleep after the max number of concurrent connections has been reached. Default to 0
 - file: Route of the file containing the urls to crawl, separated by newlines. Default to urls.txt at root of the project.
 - showResults: Show a summary of the urls with possible exposed .git or .env uris.
 
 ## Disclaimer
 
 I do not endorse or encourage the use of this crawler to engage in illegal pentesting. Before using this crawler make sure you got the proper written permission from the website's owner and make sure to consult with your lawyer.
