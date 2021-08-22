package download

import (
	"net/http"
	"strconv"
)

func stringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gopher Download Manager v1.0")
	return req, nil
}
