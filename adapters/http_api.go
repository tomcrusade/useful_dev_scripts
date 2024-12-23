package adapters

import (
	"bytes"
	"dev_scripts/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpMethod string

const (
	HttpMethodGet    HttpMethod = "GET"
	HttpMethodPost   HttpMethod = "POST"
	HttpMethodPut    HttpMethod = "PUT"
	HttpMethodDelete HttpMethod = "DELETE"
)

type HttpCall[ResultMap interface{}] struct {
	Err     error
	method  HttpMethod
	request *http.Request
}

type CallApiArgs struct {
	FullPath      string
	Token         string
	Method        HttpMethod
	RequestParams interface{}
}

func CallApi[ResponseMap interface{}](args CallApiArgs) (ResponseMap, error) {
	call := newHttpCall[ResponseMap](args.Method, args.FullPath, args.RequestParams).
		setHeader("Authorization", fmt.Sprintf("Bearer %s", args.Token)).
		setHeader("Content-Type", "application/json")
	return call.execute().parseResponse()
}

func newHttpCall[ResultMap interface{}](method HttpMethod, url string, body interface{}) HttpCall[ResultMap] {
	marshalledJsonData, err1 := json.Marshal(body)
	if err1 != nil {
		return HttpCall[ResultMap]{Err: err1}
	}

	request, err := http.NewRequest(string(method), url, bytes.NewBuffer(marshalledJsonData))
	if err != nil {
		return HttpCall[ResultMap]{Err: err}
	}

	output := HttpCall[ResultMap]{method: method, request: request}
	return output
}

func (s HttpCall[R]) setHeader(key string, value string) HttpCall[R] {
	s.request.Header.Set(key, value)
	return s
}

func (s HttpCall[R]) execute() httpExecuteResult[R] {
	var client = &http.Client{}
	response, err2 := client.Do(s.request)
	return httpExecuteResult[R]{Err: err2, method: s.method, response: response}
}

type httpExecuteResult[ResultMap interface{}] struct {
	Err      error
	method   HttpMethod
	response *http.Response
}

func (r httpExecuteResult[R]) parseResponse() (R, error) {
	var result R
	if r.Err != nil {
		return result, r.Err
	}

	if r.response.Body != nil {
		defer r.response.Body.Close()
		err1 := json.NewDecoder(r.response.Body).Decode(&result)
		if err1 != nil && err1 != io.EOF {
			return result, err1
		}
	}
	if r.response.StatusCode < 200 || r.response.StatusCode >= 300 {
		return result, fmt.Errorf(
			"%s api returns status code %v with parsedResult %v",
			r.method,
			r.response.StatusCode,
			entity.ConvertToJSON(result),
		)
	}

	return result, nil
}
