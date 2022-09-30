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

		// TODO determine if io.Copy is actually a lazy operation
		// or if it will blow our memory up
		// _, _ = io.Copy(gzipWriter, input)

		buffer := make([]byte, 1024)
		for {
			n, err := input.Read(buffer)
			_, _ = gzipWriter.Write(buffer[:n])
			if err == io.EOF {
				return
			} else if err != nil {
				return
			}
		}
	}()
	return pipeReader
}
