# Bencode

A fast and secure Go library to encode and parse [Bencode](https://en.wikipedia.org/wiki/Bencode).

Inspired by the great work done by [@jackpal/bencode-go](https://github.com/jackpal/bencode-go/) and [@marksamman/bencode](https://github.com/marksamman/bencode).

## Examples

Here are some usage examples:

```go
import "github.com/stefanovazzocell/bencode"

// Parse a string
bencode.NewParserFromString("10:helloworld").AsString() // "helloworld", nil
// Parse an object
bencode.NewParserFromString("d4:name5:Alice3:agei35ee").AsInterface() // map[string]{}{ "name": "Alice", "age": 35 }, nil

// Parse a map from an io.Reader
// fileReader implements io.Reader and returns "li1ei2ei3ee"
bencode.NewParserFromReader(fileReader).AsList() // []interface{1, 2, 3}


// Encode an object
encoder, err := bencode.NewEncoderFromInterface([]interface{}{1,2,3})
if err != nil {
    // TODO: Handle error
}
encoder.String() // "li1ei2ei3ee"
```

## Aims

As per the introduction I aim to make this library *fast* and *secure*.
Here I will address how I plan to achieve those goals.

### Secure

Since this library contains a parser that might be used to read user-generated content it's important for it to be well tested.

- Near-100% code coverage (`make test`, coverage: up_to* 96.1%).
- Extensively fuzzed (`make fuzz`, coverage: 99.3%).
- Checked for security issues with [`gosec`](https://github.com/securego/gosec) (`make security`).

*up_to: one of the tests uses an intentionally unreliable `io.Reader`.

### Performance

This library needs to perform well as it might need to encode/decode a large amount of data efficiently. Benchmarks are available with `make bench`.

```
$ make bench
go test -run=^$ -cover -bench .
goos: linux
goarch: amd64
pkg: github.com/stefanovazzocell/bencode
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkEncoder/torrentString-12               13454383        88.96 ns/op
BenchmarkEncoder/complexMap-12                   1384812       867.4 ns/op
BenchmarkReaderParser/complexMap-12               545340      1873 ns/op
BenchmarkReaderParser/torrentString-12           3972987       329.2 ns/op
BenchmarkReaderParser/torrentStringAsString-12   4413070       269.6 ns/op
BenchmarkStringParser/complexMap-12               730975      1386 ns/op
BenchmarkStringParser/torrentString-12          14311644        77.85 ns/op
BenchmarkStringParser/torrentStringAsString-12  26242033        42.01 ns/op
```

Furthermore you can profile specific components with the following:
- `make profileEncoder`
- `make profileStringParser`

## Caveats

There are some things to consider
- This library can encode `int` and `uint` as well as all their variations (i.e.: `int64`, `uint16`, ...) but it can only parse numbers of type `int`.
- This library can only encode maps of type `map[string]interface{}` and slices of type `[]interface{}`.
- Additional data after the initial parse will be ignored, unless another parse operation (such as `.AsList()`) is called on the same parser.
