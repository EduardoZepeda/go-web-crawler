# Go crawler for .env and .git

This crawler was inspired by [scanning-26-million-domains-for-exposed-env-files](https://hackernoon.com/scanning-26-million-domains-for-exposed-env-files) article. It uses concurrency to crawl a list of domains and check for exposed .env and .git uris, in plain or www subdomains.

## Quickstart

### Clone the project or build an executable

As any go project you can run it directly or compile it to produce a binary file.

 ```bash
git clone https://github.com/eduardoZepeda/go-web-crawler
cd go-web-crawler/
 ```

### You need a file with urls

You will need a file with domain names separated by newlines. The default expected file name is *urls.txt*, located at the root of the project. Please see flag options to specify another file name.

For instance, consider a file named *urls.txt* at the root of the project or executable with the following domains:

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

 You can also execute the previously generated binary file

 ### Getting the results

 A list of all vulnerable sites will be printed at the end of the code execution. If nothing was found, no output will be shown.
 
 ## Acceptance criteria
 
 The crawler consider a successful response any response with a code status between 200 and 300. The crawler ignores any redirection.

 Note: I'm aware this could lead to some false positives. Please modify the code according to your own needs.
 
 ## Flag options
 
 The following options can be used to customize the crawler behaviour.
 
 - logLevel: The log Level. Valid values: from 1 to 6, ascending verbosity
 - concurrent: Max number of concurrent requests. Default to 10
 - reqTimeout: Timeout (in seconds) before http request is aborted. Default to 5
 - connTimeout: Timeout (in seconds) before opening a new http connection. Default to 10
 - sleep: Timeout (in seconds) to sleep after the max number of concurrent connections has been reached. Default to 0
 - file: Route of the file containing the urls to crawl, separated by newlines. Default to *urls.txt* at root of the project.
 - showResults: Show a summary of the urls with possible exposed .git or .env uris at the end of the execution process.

 ## Quick start with Docker

 To run this project using docker just follow the usual procedure. Make sure the urls file exists at the root of the project.

 ```bash
 docker build . -t <name>
 ``` 

 After that you can run the crawler. You're free to pass any number of flags and its corresponding values.

 ```bash
 docker run --rm <name> [<flag>=<value>...]
 ``` 

 ## Disclaimer
 
 I do not endorse or encourage the use of this crawler to engage in illegal pentesting. Before using this crawler make sure you got the proper written permission from the website's owner and make sure to consult with your lawyer.
