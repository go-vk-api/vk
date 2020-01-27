package vk

import (
	"testing"
)

func TestRequestParams_UrlValues(t *testing.T) {
	cases := []struct {
		in   RequestParams
		want string
	}{
		{
			RequestParams{
				"boolean": true,
				"int":     108,
				"string":  "4 8 15 16 23 42",
			},
			"boolean=true&int=108&string=4+8+15+16+23+42",
		},
	}

	for _, c := range cases {
		urlValues, err := c.in.UrlValues()

		if err != nil {
			t.Error(err)
		}

		if len(urlValues) != len(c.in) {
			t.Errorf("len(urlValues) == %d, want %d", len(urlValues), len(c.in))
		}

		if urlValues.Encode() != c.want {
			t.Errorf("urlValues.Encode() == %q, want %q", urlValues.Encode(), c.want)
		}
	}
}
