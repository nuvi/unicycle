package fetch

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/nuvi/unicycle/error_ext"
)

type FetchOptions struct {
	Method            string
	Query             map[string]string
	Headers           map[string]string
	Body              io.Reader
	Timeout           *time.Duration
	Logging           bool
	AcceptBadResponse bool
}

func LogResponseDetails(response *http.Response) {
	if response == nil {
		log.Println("LogResponseDetails() error: response is nil")
	} else {
		if response.Request == nil {
			log.Println("LogResponseDetails() error: response.Request is nil")
		} else {
			if response.Request.URL == nil {
				log.Println("LogResponseDetails() error: response.Request.URL is nil")
			} else {
				log.Println(response.Request.URL)
			}
		}
		log.Println(response.Status)
		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("LogResponseDetails() error: could not read body:", err)
		} else {
			log.Println(string(responseBodyBytes))
		}
	}
}

func ResponseOk(response *http.Response) (bool, error) {
	if response == nil {
		return false, errFetchNilResponse
	} else {
		if (response.StatusCode < http.StatusOK) || (http.StatusMultipleChoices <= response.StatusCode) {
			return false, newFetchError(BadResponseError{StatusCode: response.StatusCode}, response)
		}
		return true, nil
	}
}

func LogPossibleFetchError(err error) bool {
	log.Println(err)
	if fetchError := error_ext.ErrorAs[FetchError](err); fetchError != nil {
		LogResponseDetails(fetchError.Response)
		return true
	}
	return false
}

func AppendQueryParams(rawUrl string, queryParams map[string]string) (string, error) {
	trueUrl, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl, err
	}

	if len(queryParams) > 0 {
		query := trueUrl.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		trueUrl.RawQuery = query.Encode()
	}

	return trueUrl.String(), nil
}

// Fetch simplifies common http requests and associated error checking
func Fetch(rawUrl string, options FetchOptions) (*http.Response, error) {
	trueUrl, err := AppendQueryParams(rawUrl, options.Query)
	if err != nil {
		return nil, err
	}

	if options.Method == "" {
		options.Method = "GET"
	}

	request, err := http.NewRequest(options.Method, trueUrl, options.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range options.Headers {
		request.Header.Add(key, value)
	}

	timeout := time.Minute
	if options.Timeout != nil {
		timeout = *options.Timeout
	}

	client := http.Client{
		Timeout: timeout,
	}
	if options.Logging {
		log.Println(options.Method, trueUrl)
	}
	response, err := client.Do(request)
	if err != nil {
		return response, newFetchError(err, response)
	}

	if options.Logging {
		log.Println(response.Status)
	}

	return response, nil
}
