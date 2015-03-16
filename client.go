package loghttp

import (
	"fmt"
	"net/http"
	"os"
)

// LogRoundTripper logs all HTTP requests and responses before delegating them to an underlying
// RoundTripper implementation.
type LogRoundTripper struct {
	inner   http.RoundTripper
	enabled bool
}

// NewLogRoundTripper returns a LogRoundTripper that delegates to the default transport. It begins
// enabled unless the LOGHTTP_DISABLED environment variable is set to a nonempty value.
func NewLogRoundTripper() *LogRoundTripper {
	return &LogRoundTripper{
		inner:   http.DefaultTransport,
		enabled: os.Getenv("LOGHTTP_DISABLED") == "",
	}
}

// RoundTrip executes a single HTTP request and logs it using current configuration.
func (rt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error

	if rt.enabled {
		fmt.Printf("%s %s %s\n", request.Method, request.URL.RequestURI(), request.Proto)
		request.Header.Write(os.Stdout)

		if request.Body != nil {
			fmt.Println()
			request.Body, err = outputBody(request.Body, request.Header)
			if err != nil {
				return nil, err
			}
		}
	}

	response, err := rt.inner.RoundTrip(request)
	if err != nil {
		return response, err
	}

	if rt.enabled {
		fmt.Printf("%s %s\n", response.Proto, response.Status)
		response.Header.Write(os.Stdout)
		response.Body, err = outputBody(response.Body, response.Header)
	}

	return response, err
}

// NewClient returns an HTTP client that's initialized to use the LogRoundTripper. All other
// settings are left as the default.
func NewClient() http.Client {
	return http.Client{
		Transport: NewLogRoundTripper(),
	}
}
