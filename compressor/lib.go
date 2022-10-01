package compressor

import (
	"compress/gzip"
	"io"
)

func AsReader(input io.ReadCloser) io.Reader {
	// There is no middleware pattern for readers, but A Pipe
	// can be used to effectively chain two readers together.
	// Everything passed to the write end of the pipe is accessible
	// to the read end of the pipe.
	pipeReader, pipeWriter := io.Pipe()
	// The gzip writer will write data into the pipe, which will
	// be accessible any downstream consumer of a reader.
	gzipWriter := gzip.NewWriter(pipeWriter)
	go func() {
		// Pipes are synchronous. You will find if you do not have
		// a goroutine processing them, you will have mysteriously
		// blocking code.
		defer input.Close()
		defer pipeWriter.Close()
		defer gzipWriter.Close()
		// One can manually buffer the input reader to the writer,
		// however, the standard library's version checks the underlying
		// types of the Readers and Writers for optimized method calls
		if _, err := io.Copy(gzipWriter, input); err != nil {
			panic(err)
		}
	}()
	return pipeReader
}
