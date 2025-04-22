package request

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLine(t *testing.T) {
	// TEST: Valid Request Line
	rq := RequestLine{}
	reader := &chunkReader{
		data:            "VXP/1.0 OBTAIN www.google.com",
		numBytesPerRead: 10,
	}
	r, rq, err := rq.ParseRequestLine([]byte(reader.data))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "VXP/1.0", rq.Version)
	assert.Equal(t, "OBTAIN", rq.Method)
	assert.Equal(t, "www.google.com", rq.Url)

	// TEST: Valid Request Line
	rq = RequestLine{}
	reader = &chunkReader{
		data:            "VXP/1.0 BYE localhost:42069",
		numBytesPerRead: 10,
	}
	r, rq, err = rq.ParseRequestLine([]byte(reader.data))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "VXP/1.0", rq.Version)
	assert.Equal(t, "BYE", rq.Method)
	assert.Equal(t, "localhost:42069", rq.Url)

	// TEST: Invalid Method
	rq = RequestLine{}
	reader = &chunkReader{
		data:            "VXP/1.0 GET localhost:42069",
		numBytesPerRead: 10,
	}
	r, rq, err = rq.ParseRequestLine([]byte(reader.data))
	fmt.Printf("Request line should fail %s", rq.Method)
	require.Error(t, err)

	// TEST: Invalid Version
	rq = RequestLine{}
	reader = &chunkReader{
		data:            "VXP/1.1 GET localhost:42069",
		numBytesPerRead: 10,
	}
	r, rq, err = rq.ParseRequestLine([]byte(reader.data))
	require.Error(t, err)
	require.NotNil(t, r)

	// TEST: Invalid Length
	rq = RequestLine{}
	reader = &chunkReader{
		data:            "VXP/1.0 GET localhost:42069 more data",
		numBytesPerRead: 10,
	}
	r, rq, err = rq.ParseRequestLine([]byte(reader.data))
	require.Error(t, err)
	require.NotNil(t, r)

}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}
