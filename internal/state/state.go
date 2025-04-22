package state

import (
	"errors"

	"github.com/hconn7/vxp/internal/headers"
	"github.com/hconn7/vxp/internal/request"
	"github.com/hconn7/vxp/types/constants"
)

type Request struct {
	state       constants.StateId
	requestLine request.RequestLine
	headers     headers.Headers
	Body        []byte
}

func (r *Request) ParseState(data []byte) (int, error, Request) {
	switch r.state {
	case constants.STATE_REQLINE:
		n, reqLine, err := r.requestLine.ParseRequestLine(data)
		if err != nil {
			return 0, errors.New("Error parsing reqLine"), Request{}

		}
	}
}
