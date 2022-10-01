package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzCompressorAsReader(f *testing.F) {
	f.Add("Hi, I'm Paul")
	f.Fuzz(func(t *testing.T, fuzzString string) {
		// Simulate generic ReaderCloser
		var input io.ReadCloser = ioutil.NopCloser(strings.NewReader(fuzzString))
		// Convert this reader to a CompressionReader for downstream consumers
		reader := AsReader(input)
		// Validate that the compression happened by decompressing and comparing with the input
		actualOutput, err := io.ReadAll(reader)
		if err != nil {
			t.Error(err)
			return
		}
		gzipReader, _ := gzip.NewReader(bytes.NewReader(actualOutput))
		actualBytes, _ := io.ReadAll(gzipReader)
		assert.Equal(t, []byte(fuzzString), actualBytes)
	})
}
