package headers

import (
	_ "fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewHeaders() Headers {
	return make(Headers)
}

func TestRequestLineParse(t *testing.T) {
	// Test: Valid single heade
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, _, err = headers.Parse(data)
	require.NoError(t, err)
	_, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.True(t, done)

	data = []byte("Host:localhost:42069\r\nContent:foo\r\nSexy:beast\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "foo", headers["content"])
	assert.Equal(t, 13, n)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 12, n)
	assert.Equal(t, "beast", headers["sexy"])

	data = []byte("Hostlocalhost_42069\r\n\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 21, n)
	assert.False(t, done)
	_, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.True(t, done)

	data = []byte("H©st: localhost:42069\r\n\r\n")
	targetLen := len(data) - 2
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, targetLen, n)
	data = data[n:]
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)

	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42420\r\n\r\n")
	n, _, _ = headers.Parse(data)
	data = data[n:]
	_, _, _ = headers.Parse(data)
	assert.Equal(t, "localhost:42069, localhost:42420", headers["host"])

}
