package request

import (
	"errors"
	"fmt"
	"github.com/posixenjoyer/httpfromtcp/internal/headers"
	_ "github.com/stretchr/testify/assert"
	"io"
	"strings"
	_ "testing"
)

func (r RequestLine) IsEmpty() bool {
	return r == RequestLine{}
}

func (r *Request) parse(message []byte) (int, error) {
	var count int
	var err error
	var parsedRequestLine RequestLine

	if r.State == DONE {
		return 0, errors.New("trying to read in DONE state")
	}

	if r.State == INITIALIZED {
		parsedRequestLine, count, err = parseRequestLine(message)

		if !parsedRequestLine.IsEmpty() {
			r.State = DONE
			r.RequestLine = parsedRequestLine
		}
	} else {
		return 0, errors.New("unknown state")
	}

	return count, err
}

func (r *Request) checkEOFState(buffer []byte, n int) (int, error) {
	if n > 0 {
		requestLine, count, err := parseRequestLine(buffer)
		if requestLine.IsEmpty() {
			return count, err
		}

		r.State = DONE
		r.RequestLine = requestLine
		return count, nil
	}

	return 0, errors.New("connection terminated prematurely")
}

func RequestFromReader(request io.Reader) (*Request, error) {
	lastCount := 0
	buffer := make([]byte, BUF_SZ)
	readToIndex := 0
	parsedRequest := Request{
		RequestLine: RequestLine{},
		Body:        make([]byte, 256),
		Headers:     make(map[string]string),
		State:       INITIALIZED,
	}

	for parsedRequest.State == INITIALIZED {
		bufLen := len(buffer)
		if readToIndex == bufLen {
			newBuffer := make([]byte, bufLen*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		lastCount, err := request.Read(buffer[readToIndex:])
		readToIndex += lastCount

		if err != nil {
			if err == io.EOF {
				parsedRequest.State = CHECK
				break
			}
		}
		count, err := parsedRequest.parse(buffer)
		if err != nil {
			return &Request{}, err
		}

		if count > 0 {
			parsedRequest.State = DONE
			copy(buffer, buffer[count:readToIndex])
			readToIndex -= count
		}
	}

	if parsedRequest.State == CHECK {
		_, err := parsedRequest.checkEOFState(buffer, lastCount)
		if err != nil {
			fmt.Printf("EOF read parse failed: %+v\n", err)
			return &Request{}, err
		}
	}
	return &parsedRequest, nil
}

func parseRequestLine(requestLine []byte) (RequestLine, int, error) {

	if !strings.Contains(string(requestLine), "\r\n") {
		return RequestLine{}, 0, nil
	}

	requestStr := strings.Split(string(requestLine), "\r\n")[0]
	requestStr = strings.Trim(requestStr, "\r\n")
	requestArgs := strings.Split(string(requestStr), " ")

	if len(requestArgs) < 3 {
		err := fmt.Errorf("error: Not enough fields in request-line")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	if len(requestArgs) > 3 {
		err := fmt.Errorf("error: Too many arguments for request-line")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	if !verifyMethod(requestArgs[0]) {
		err := fmt.Errorf("error: Method name verification failed")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	version := strings.TrimRight(requestArgs[2], "\r\n ")

	if !verifyVersion(version) {
		err := fmt.Errorf("error: Invalid HTTP Version: %s", requestArgs[2])
		fmt.Printf("RS: %s, Version: %s\nLen(version): %d\n", requestStr, version, len(version))
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}
	parsedReqestLine := RequestLine{
		Method:        requestArgs[0],
		RequestTarget: requestArgs[1],
		HttpVersion:   strings.Split(version, "/")[1],
	}

	return parsedReqestLine, len(requestStr) + len("\r\n"), nil
}
