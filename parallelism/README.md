# Parallelism

## Usage

First install the binary. With the asumption that one's shell has
the `${GOPATH}/bin` directory in one's path, you may execute the
built CLI application and pass CSV filenames as arguments to parse.

> Note, you can temporarily add this directory to your path if
> preferable:
> `export PATH="$(go env GOPATH)/bin:${PATH}"`

```bash
$ go install .
$ parallelism \
  ./testdata/TestProcessEligibleChannel2_0_TestProcessEligibleChannel2_0_CODES.csv \
  ./testdata/TestProcessEligibleChannel2_4_TestProcessEligibleChannel2_4_CODES.csv
```

Alternatively, one can simply "go run" the source-code instead

```bash
$ go run .  \
  ./testdata/TestProcessEligibleChannel2_0_TestProcessEligibleChannel2_0_CODES.csv \
  ./testdata/TestProcessEligibleChannel2_4_TestProcessEligibleChannel2_4_CODES.csv
```

## Tests

Unit tests are executed with

```bash
$ go test -v --count=1 .
```

They iterate over the canned test-fixtures, located in [./testdata/](./testdata), that include both
unique data and duplicate data. They assert errors should be returned
when duplicates are present, and conversly assert when they shouldn't.
