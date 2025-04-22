package request

import (
	"errors"
	"fmt"
	"strings"
)

type RequestLine struct {
	Version string
	Method  string
	Url     string
}

// VXP/1.0 METHOD URL
func (rq *RequestLine) ParseRequestLine(data []byte) (parsed int, requestLine RequestLine, err error) {
	s := strings.TrimSpace(string(data))
	dat := strings.Split(s, " ")
	if len(dat) != 3 {
		fmt.Printf("Malformed Request Line With lenght: %v\n", len(dat))
		return 0, RequestLine{}, errors.New("Data is malformed: Length")

	}

	version := dat[0]
	versCheck := strings.Split(version, "/")
	if versCheck[0] != "VXP" {
		return 0, RequestLine{}, fmt.Errorf("Version header is malformed\n Wanted VXP Got: %s\n", versCheck[0])

	}
	if len(versCheck) != 2 || versCheck[0] != "VXP" || versCheck[1] != "1.0" {
		return 0, RequestLine{}, fmt.Errorf("invalid version format: got %s", version)
	}

	rq.Version = version
	rq.Method = dat[1]
	rq.Url = dat[2]
	requestLine = *rq

	valid := []string{"PING", "CAST", "OBTAIN", "VALID", "OMIT", "BYE"}
	var v bool
	for _, val := range valid {
		if rq.Method == val {
			v = true
		}
	}
	if !v {
		fmt.Println("Request Method doesn't exist")
		return 0, RequestLine{}, errors.New("Invalid Req Method")
	}
	return len(data), requestLine, nil

}
