package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		log.Fatalf("No files were provided for processing")
	} else if err := ProcessFiles(context.Background(), files...); err != nil {
		log.Fatalf("Error while processing files: %+v", err)
	}
	log.Printf("No duplicates were found in: %s", strings.Join(files, ", "))
}

type ProcessedCode struct {
	Code string
	Err  error
}

func ProcessFiles(ctx context.Context, filenames ...string) error {
	// Create a cancelable context from the one we are provided. If the caller wishes
	// to cancel we should act upon that, but we must be able to terminate early ourselves
	// as well.
	cancelableCtx, cancelParentCtx := context.WithCancel(ctx)
	defer cancelParentCtx()
	group, groupCtx := errgroup.WithContext(cancelableCtx)
	// Make a single channel for consuming codes. This will allow our workers to "fan-in"
	codeChannel := make(chan string)
	// For every csv-file we have been asked to process, we iterate and open them.
	// We create a dedicated supplying channel for each. These suppliers are tied
	// to our consuming channel
	for _, filename := range filenames {
		filename := filename
		group.Go(func() error {
			file, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()
			reader := csv.NewReader(file)
			supplyingChannel := csvReaderAsSupplier(groupCtx, reader)
			for code := range supplyingChannel {
				if code.Err != nil {
					return code.Err
				}
				codeChannel <- code.Code
			}
			return nil
		})
	}

	// WaitGroups, and ErrGroups as well, do not operate with channels. They have blocking
	// `Wait` methods. This means that we need to create a channel we manage so we can mux
	// both consuming codes and termination cases
	waitGroupChannel := make(chan error)
	go func() {
		defer close(codeChannel)
		defer close(waitGroupChannel)
		defer cancelParentCtx()
		waitGroupChannel <- group.Wait()
	}()

	// We use a map[string]struct{} as the communal standard for a Set in Go.
	codeSet := make(map[string]struct{})
	for {
		select {
		case currentCode := <-codeChannel:
			if _, isPresent := codeSet[currentCode]; isPresent {
				return fmt.Errorf("duplicate code: %s", currentCode)
			}
			codeSet[currentCode] = struct{}{}
		// If the wait group is done, we can bubble any error that occurred
		// because it will either be related to the file processing OR no error
		// will be present indicating we've completed our work
		case err := <-waitGroupChannel:
			return err
		}
	}

}

func csvReaderAsSupplier(ctx context.Context, reader *csv.Reader) chan ProcessedCode {
	codeChannel := make(chan ProcessedCode)
	go func() {
		defer close(codeChannel)
		// Consume the CSV header while also determining which column contains the value
		// we want to process
		headerLine, err := reader.Read()
		if err != nil {
			codeChannel <- ProcessedCode{
				Err: err,
			}
			return
		}
		codeIndex := indexOf("code", headerLine)
		if codeIndex < 0 {
			codeChannel <- ProcessedCode{
				Err: fmt.Errorf("csv header did not denote a `code` field"),
			}
			return
		}
		for {
			// Iterate through all records and pass them to the consuming channel
			currentRecord, err := reader.Read()
			if err == io.EOF {
				return
			} else if err != nil {
				codeChannel <- ProcessedCode{
					Err: err,
				}
				return
			}
			// Check if our context has been canceled and terminate if it has
			select {
			case <-ctx.Done():
				return
			default:
				codeChannel <- ProcessedCode{
					Code: currentRecord[codeIndex],
				}
			}
		}
	}()
	return codeChannel
}

func indexOf[T comparable](item T, slice []T) int {
	for i := 0; i < len(slice); i++ {
		if item == slice[i] {
			return i
		}
	}
	return -1
}
