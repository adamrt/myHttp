package request

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type RequestLine struct {
	Version string
	Method  string
	Target  string
}

// VXP/1.0 METHOD URL

func ParseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return &RequestLine{}, 0, nil
	}
	s := data[:idx]
	stringDat := string(s)
	fmt.Printf("Data string: %s", stringDat)
	requestSplit := strings.Split(stringDat, " ")
	if len(requestSplit) != 3 {
		return nil, 0, errors.New("Malformed format")
	}
	reqMethod := strings.TrimSpace(requestSplit[1])
	reqTarget := requestSplit[2]
	reqVersion := requestSplit[0]
	fmt.Printf("Version: %s\n Method: %s\n Target: %s\n", reqVersion, reqMethod, reqTarget)
	versionNum := strings.Split(reqVersion, "/")
	if versionNum[1] != "1.0" {
		return &RequestLine{}, 0, errors.New("Version is not 1.1")
	}

	validMethods := []string{"PING", "CAST", "OBTAIN", "VALID", "OMIT", "BYE"}
	valid := false
	for _, v := range validMethods {
		if reqMethod == v {
			valid = true
			break
		}
	}
	if !valid {
		return &RequestLine{}, 0, errors.New("Method is invalid")

	}

	return &RequestLine{
		Version: reqVersion,
		Target:  reqTarget,
		Method:  reqMethod}, idx + 2, nil
}
