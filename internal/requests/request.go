package request

import (
	"fmt"
	_ "github.com/stretchr/testify/assert"
	"io"
	"strings"
	_ "testing"
)

func RequestFromReader(request io.Reader) (*Request, error) {
	var parsedRequest Request

	lines, err := io.ReadAll(request)

	linesStr := string(lines)
	if err != nil {
		fmt.Printf("Request Read failed: %v\n", err)
	}
	parsedResult, err := parseRequestLine(strings.Split(linesStr, "\r\n")[0])
	if err != nil {
		fmt.Printf("Parsing failed: %v\n", err)
		parsedRequest.RequestLine = RequestLine{}
		return &parsedRequest, err
	}

	parsedRequest.RequestLine = parsedResult
	parsedRequest.Headers = make(map[string]string)
	parsedRequest.Body = make([]byte, 8)

	return &parsedRequest, nil
}

func parseRequestLine(requestLine string) (RequestLine, error) {

	requestArgs := strings.Split(requestLine, " ")

	if len(requestArgs) < 3 {
		err := fmt.Errorf("error: Not enough fields in request-line")
		return RequestLine{}, err
	}

	if len(requestArgs) > 3 {
		err := fmt.Errorf("error: Too many arguments for request-line")
		return RequestLine{}, err
	}

	fmt.Printf("requestArgs[0]: %s\n", requestArgs[0])

	if !verifyMethod(requestArgs[0]) {
		err := fmt.Errorf("error: Method name verification failed")
		return RequestLine{}, err
	}

	if !verifyVersion(requestArgs[2]) {
		err := fmt.Errorf("error: Invalid HTTP Version: %s", requestArgs[2])
		return RequestLine{}, err
	}
	parsedReqestLine := RequestLine{
		Method:      requestArgs[0],
		Target:      requestArgs[1],
		HttpVersion: strings.Split(requestArgs[2], "/")[1],
	}

	return parsedReqestLine, nil
}
