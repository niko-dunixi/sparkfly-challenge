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

func FuzzCompressor(f *testing.F) {
	f.Add("Hi, I'm Paul")
	f.Fuzz(func(t *testing.T, fuzzString string) {
		// Convert this reader to a CompressionReader for downstream consumers
		reader := AsReader(ioutil.NopCloser(strings.NewReader(fuzzString)))

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
