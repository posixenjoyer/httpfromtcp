package request

const DONE int = 2
const INITIALIZED int = 1
const CHECK int = 3
const BUF_SZ int = 8

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	State       int
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

var httpMethods []string = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}
