package crawler

import (
	"net/http"
)

type Request struct {
	httpReq *http.Request
	depth   uint32
}

func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{
		httpReq: httpReq,
		depth:   depth,
	}
}

func (r *Request) Get() *http.Request {
	return r.httpReq
}

func (r *Request) Depth() uint32 {
	return r.depth
}

func (r *Request) Valid() bool {
	return r.httpReq != nil && r.httpReq.URL != nil
}
