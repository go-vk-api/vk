package vk

import (
	"fmt"
	"net/url"
)

// RequestParams are the params for invoking methods.
type RequestParams map[string]interface{}

// URLValues translates the params to url.Values.
func (params RequestParams) URLValues() (url.Values, error) {
	values := url.Values{}

	for k, v := range params {
		values.Add(k, fmt.Sprint(v))
	}

	return values, nil
}
