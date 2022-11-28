package utils

import "net/url"

func FormatUrl(prefix string, uri string) (url.URL, error) {
	// According to go's documentation urls are recognized as [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	joinedUrl, err := url.JoinPath(prefix, uri, "/")
	if err != nil {
		return url.URL{}, err
	}
	u, err := url.ParseRequestURI(joinedUrl)
	if err != nil {
		return url.URL{}, err
	}
	return *u, nil
}
