package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var (
	client *resty.Client = resty.New()
)

type requestType int

const (
	get  requestType = 0
	post requestType = 1
)

type request struct {
	reqType     requestType
	endpoint    string
	pathParams  map[string]string
	queryParams map[string]string
	body        interface{}
	contentType string
}

func buildRequestUrl(baseEndpoint string, endpoint string) string {
	return fmt.Sprintf("%s%s", baseEndpoint, endpoint)
}

func doReq(baseEndpoint string, in request, output interface{}) (int, error) {
	r := client.R().
		SetPathParams(in.pathParams).
		SetQueryParams(in.queryParams)

	if in.body != nil {
		r.SetBody(in.body)
	}

	var resp *resty.Response
	var err error

	switch in.reqType {
	case get:
		resp, err = r.Get(buildRequestUrl(baseEndpoint, in.endpoint))
	case post:
		resp, err = r.Post(buildRequestUrl(baseEndpoint, in.endpoint))
	default:
		return http.StatusNotImplemented, errors.New("operation not supported")
	}

	if err != nil {
		return resp.StatusCode(), err
	}

	if output != nil {
		err = json.Unmarshal(resp.Body(), output)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode(), nil
}
