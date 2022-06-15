package client

import (
	"fmt"
	"os"
)

type Api struct {
	Host  string
	Token string
	Debug bool
}

type ApiResponse struct {
	Meta struct {
		TotalCount int `json:"total_count"`
	}
}

func NewApi() *Api {
	debug := false
	if os.Getenv("TF_LOG") == "DEBUG" {
		debug = true
	}
	return &Api{
		Host:  os.Getenv("DI_URL"),
		Token: os.Getenv("DI_TOKEN"),
		Debug: debug,
	}
}

func (o *Api) request(method, resource string) *Request {
	return &Request{
		Debug:         o.Debug,
		Url:           fmt.Sprintf("%s/%s", o.Host, resource),
		Method:        method,
		Authorization: o.Token,
	}
}

func (o *Api) NewRequest(method, resource string, body []byte, expectCode int) (*Request, error) {
	request := &Request{
		Debug:         o.Debug,
		Url:           fmt.Sprintf("%s/%s", o.Host, resource),
		Method:        method,
		Authorization: o.Token,
	}
	if body != nil {
		request.Body = body
	}
	err := request.Make()
	if err != nil {
		return nil, err
	}
	if request.Response.StatusCode != expectCode {
		return request, fmt.Errorf(
			"wrong statusCode from API: %d, expect: %d, resource [%s], response: %s",
			request.Response.StatusCode,
			expectCode,
			resource,
			string(request.ResponseBody),
		)
	}
	return request, nil
}

func (o *Api) NewRequestCreate(url string, data []byte) ([]byte, error) {
	request, err := o.NewRequest("POST", url, data, 201)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}

func (o *Api) NewRequestRead(url string) ([]byte, error) {
	request, err := o.NewRequest("GET", url, nil, 200)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}

func (o *Api) NewRequestUpdate(url string, data []byte) ([]byte, error) {
	request, err := o.NewRequest("PATCH", url, data, 200)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}

func (o *Api) NewRequestDelete(url string, data []byte) error {
	_, err := o.NewRequest("DELETE", url, data, 204)
	if err != nil {
		return err
	}
	return nil
}

func (o *Api) NewRequestResize(url string, data []byte) ([]byte, error) {
	request, err := o.NewRequest("POST", url, data, 201)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}

func (o *Api) NewRequestMove(url string, data []byte) ([]byte, error) {
	request, err := o.NewRequest("POST", url, data, 200)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}

func (o *Api) NewRequestUpScale(url string, data []byte) ([]byte, error) {
	request, err := o.NewRequest("POST", url, data, 201)
	if err != nil {
		return nil, err
	}
	return request.ResponseBody, nil
}
