package main

import (
	"context"
	"fmt"
	"testing"
)

func TestFiles(t *testing.T) {
	t.Run("Failing On Duplicates", func(t *testing.T) {
		if err := ProcessFiles(context.Background(), formatStockFixtureFilenames(0, 1)...); err == nil {
			t.Error("Should have failed because files 0 and 1 contain a duplicate")
		}
	})

	t.Run("Successful Combinations", func(t *testing.T) {
		t.Run("All but 0", func(t *testing.T) {
			if err := ProcessFiles(context.Background(), formatStockFixtureFilenames(1, 2, 3, 4)...); err != nil {
				t.Errorf("Should not have failed, but did: %+v", err)
			}
		})
		t.Run("All but 1", func(t *testing.T) {
			if err := ProcessFiles(context.Background(), formatStockFixtureFilenames(0, 2, 3, 4)...); err != nil {
				t.Errorf("Should not have failed, but did: %+v", err)
			}
		})
	})
}

func formatStockFixtureFilenames(indexes ...int) []string {
	results := make([]string, len(indexes))
	for i := 0; i < len(indexes); i++ {
		results[i] = formatStockFixture(indexes[i])
	}
	return results
}

func formatStockFixture(index int) string {
	return fmt.Sprintf(
		"./testdata/TestProcessEligibleChannel2_%d_TestProcessEligibleChannel2_%d_CODES.csv",
		index,
		index,
	)
}
