package compressor

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

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

func TestCompressorAsReader(t *testing.T) {
	t.Run("Random File", func(t *testing.T) {
		sizesInMB := []int{1, 10, 20, 100, 1000}
		for _, currentMB := range sizesInMB {
			currentMB := currentMB
			t.Run(fmt.Sprintf("%dmb", currentMB), func(t *testing.T) {
				t.Parallel()
				// Create our testable reader with an expected sha256 sum to validate
				inputFile, expectedSha256Sum := generateRandomFile(t, 1024*1024*currentMB)
				defer func() {
					inputFile.Close()
					_ = os.Remove(inputFile.Name())
				}()
				// Create our Compression Reader
				reader := AsReader(inputFile)
				// Unzip and hash the bytes, compare them against the generated hash
				// to be sure we're doing things correctly
				gzipReader, err := gzip.NewReader(reader)
				if err != nil {
					panic(err)
				}
				shaHasher := sha256.New()
				if _, err := io.Copy(shaHasher, gzipReader); err != nil {
					panic(err)
				}
				actualSha256Sum := hex.EncodeToString(shaHasher.Sum(nil))
				assert.Equal(t, expectedSha256Sum, actualSha256Sum)
			})
		}
	})
}

func generateRandomFile(t *testing.T, totalBytes int) (file *os.File, expectedSha256Sum string) {
	t.Helper()
	// Create a new random file anc create a hasher. We will use a multi-writer
	// to store the random bytes and generate a sum we can check in our test
	random := rand.New(rand.NewSource(time.Now().Unix()))
	file, err := os.CreateTemp("", fmt.Sprintf("*-random-%d.bin", totalBytes))
	if err != nil {
		panic(err)
	}
	shaHasher := sha256.New()
	writer := io.MultiWriter(file, shaHasher)
	// Create our buffer and interate while we haven't hit our byte count
	bytesWritten := 0
	byteBuffer := make([]byte, 1024)
	for bytesWritten < totalBytes {
		// Create an intermedate buffer that is dynamically sized so we can write exactly
		// the number of bytes we want.
		intermediateBuffer := byteBuffer[:min(len(byteBuffer), totalBytes-bytesWritten)]
		_, _ = random.Read(intermediateBuffer)
		n, err := writer.Write(intermediateBuffer)
		if err != nil {
			panic(err)
		}
		bytesWritten += n
	}
	// Reset the file's seek location for reading at the beginning
	_, _ = file.Seek(0, 0)
	return file, hex.EncodeToString(shaHasher.Sum(nil))
}

type Int interface {
	byte | int | int32 | int64 | uint | uint32 | uint64
}

func min[T Int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// FIXME: This may be a potentially huge optimzation of the randomness generation,
// but for now we'll leave our naive implementation because the test would still
// need to consume both readers or we end up with blocked logic.
//
// func RandomBytes(total int) (input io.ReadCloser, expected io.Reader) {
// 	random := rand.New(rand.NewSource(time.Now().Unix()))
// 	pipeReader, pipeWriter := io.Pipe()
// 	teeReader := io.TeeReader(random, pipeWriter)
// 	return ioutil.NopCloser(teeReader), pipeReader
// }
