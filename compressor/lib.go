package compressor

import (
	"compress/gzip"
	"io"
)

func AsReader(input io.ReadCloser) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	gzipWriter := gzip.NewWriter(pipeWriter)
	go func() {
		defer input.Close()
		defer pipeWriter.Close()
		defer gzipWriter.Close()
		if _, err := io.Copy(gzipWriter, input); err != nil {
			panic(err)
		}
	}()
	return pipeReader
}
