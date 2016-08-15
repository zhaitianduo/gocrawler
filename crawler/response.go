package crawler

import (
	"net/http"
)

type Response struct {
	response *http.Response
	depth    uint32
}

func NewResponse(response *http.Response, depth uint32) *Response {
	return &Response{
		response: response,
		depth:    depth,
	}
}

func (res *Response) Get() *http.Response {
	return res.response
}

func (res *Response) Depth() uint32 {
	return res.depth
}

func (res *Response) Valid() bool {
	return res.response != nil && res.response.Body != nil
}
