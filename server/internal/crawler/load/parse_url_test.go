package load_test

import (
	"net/url"
	"testing"

	"github.com/dzykatsha/go-web-crawler/internal/crawler/load"
)

func TestParseUrl_SameHost(t *testing.T) {
	// arrange
	baseUrl, _ := url.Parse("https://realpython.com/")
	urlString := "https://realpython.com/examples"

	// act
	result, err := load.ParseUrl(urlString, *baseUrl)

	// assert
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	if result != urlString {
		t.Errorf("wrong result: %v", result)
	}
}

func TestParseUrl_DifferentHost(t *testing.T) {
	// arrange
	baseUrl, _ := url.Parse("https://realpython.com/")
	urlString := "https://msdn.com/examples"

	// act
	_, err := load.ParseUrl(urlString, *baseUrl)

	// assert
	if err == nil {
		t.Errorf("error must not be nil: %v", err)
	}

	if err.Error() != load.HostMismatchError {
		t.Errorf("error must be WrongHostError, but: %v", err)
	}
}

func TestParseUrl_NoHost(t *testing.T) {
	// arrange
	baseUrl, _ := url.Parse("https://realpython.com/")
	urlString := "/examples"

	// act
	result, err := load.ParseUrl(urlString, *baseUrl)

	// assert
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	if result != "https://realpython.com/examples" {
		t.Errorf("wrong result: %v", result)
	}
}
