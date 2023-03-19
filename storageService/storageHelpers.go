package storageService_helpers

import (
	"bytes"
	"net/http"
)

type QueryParams struct {
	QueryParamName string
	QueryParamVal  string
}

func DoRequest(method, path string, queryParams []QueryParams, body *bytes.Buffer) (*http.Request, error) {
	if body == nil {
		body = nil
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return &http.Request{}, err
	}

	for _, v := range queryParams {
		q := req.URL.Query()
		q.Add(v.QueryParamName, v.QueryParamVal)
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}
