package request

import (
	"errors"
	"fmt"
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

		if r.RequestLine.IsEmpty() {
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
		if readToIndex == BUF_SZ {
			buffer = append(buffer, make([]byte, readToIndex)...)
		}

		lastCount, err := request.Read(buffer[readToIndex:])
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
		}
		readToIndex += lastCount
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

	requestStr := string(requestLine)

	if !strings.Contains(requestStr, "\r\n") {
		return RequestLine{}, 0, nil
	}
	requestArgs := strings.Split(string(requestLine), " ")

	if len(requestArgs) < 3 {
		err := fmt.Errorf("error: Not enough fields in request-line")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	if len(requestArgs) > 3 {
		err := fmt.Errorf("error: Too many arguments for request-line")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	fmt.Printf("requestArgs[0]: %s\n", requestArgs[0])

	if !verifyMethod(requestArgs[0]) {
		err := fmt.Errorf("error: Method name verification failed")
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}

	if !verifyVersion(requestArgs[2]) {
		err := fmt.Errorf("error: Invalid HTTP Version: %s", requestArgs[2])
		return RequestLine{}, len(requestStr) + len("\r\n"), err
	}
	parsedReqestLine := RequestLine{
		Method:        requestArgs[0],
		RequestTarget: requestArgs[1],
		HttpVersion:   strings.Split(requestArgs[2], "/")[1],
	}

	return parsedReqestLine, len(requestStr) + len("\r\n"), nil
}
