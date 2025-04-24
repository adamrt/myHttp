package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/hconn7/vxp/internal/headers"
	"github.com/hconn7/vxp/types/constants"
)

type Request struct {
	State          constants.StateId
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	BodyLengthRead int
}

func (r *Request) ParseState(data []byte) (int, error) {
	switch r.State {
	case constants.STATE_INIT:
		reqLine, n, err := ParseRequestLine(data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.State = constants.STATE_HEADERS
		return n, nil
	case constants.STATE_HEADERS:
		n, done, err := r.Headers.ParseHeaders(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = constants.STATE_BODY
		}
		return n, nil
	case constants.STATE_BODY:
		contentLenStr, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.State = StateDone
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(contentLenStr)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.BodyLengthRead += len(data)
		if r.BodyLengthRead > contentLen {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.BodyLengthRead == contentLen {
			r.State = StateDone
		}
		return len(data), nil

	}
	return 0, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != constants.STATE_DONE {
		n, err := r.ParseState(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}

	}
	return totalBytesParsed, nil
}
func HandleConn(reader io.Reader) (*Request, error) {
	fmt.Println("conn accepted")
	buf := make([]byte, 8)
	readToIndex := 0

	req := Request{
		State:       constants.STATE_INIT,
		Headers:     headers.NewHeaders(),
		RequestLine: RequestLine{},
	}

	for req.State != constants.STATE_DONE {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != constants.STATE_DONE {
					return nil, errors.New("state not done")
				}

			}
			return &Request{}, err
		}
		readToIndex += n
		parsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			fmt.Printf("Bad request: %v\n", err)
			return nil, nil
		}
		if parsed > 0 {
			copy(buf, buf[parsed:readToIndex])
			readToIndex -= parsed
		}
	}
	// TODO: Insert a router method to grab request line
	return &req, nil
}
