# Compressor

## Usage

Import this package and use it to convert an uncompressed reader
and it will zip the bytes as they are read

```golang
package main

import (
  "github.com/paul-nelson-baker/sparkfly-challenge/compressor"
)

func main() {
  yourUncompressedReader := // this will be a io.ReadCloser of your choosing
  compressedReader := compressor.AsReader(yourUncompressedReader)
  _ = consumeReader(compressedReader)
}

func consumeReader(r io.Reader) error {
  // ...
}
```

## Tests

There are unit tests and a fuzz test in this location. To execute the fuzz test
execute from your terminal:

```bash
$ go test -fuzz .
```

To execute the unit tests:

```bash
$ go test -v -count=1
```

> Note, in an effort to conserve RAM in exchange for hard-disk, large temp
> files are created to mimic the conditions under which this function would
> be used. While there is a cleanup code within the test-logic, there is
> no guarantee the function will still be executed should the test be externally
> terminated early.
