package httputil

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Post issues a POST to the specified URL.
func Post(rd RequestDoer, url string, params url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := rd.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
