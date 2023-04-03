package load

import (
	"errors"
	"net/url"
)

const HostMismatchError = "host mismatch"

func ParseUrl(urlString string, baseUrl url.URL) (string, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}

	if parsedUrl.Host == "" {
		parsedUrl.Host = baseUrl.Host
		parsedUrl.Scheme = baseUrl.Scheme
	}

	if parsedUrl.Host != baseUrl.Host {
		return urlString, errors.New(HostMismatchError)
	}

	return parsedUrl.String(), nil
}
