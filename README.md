# Bencode

A fast and secure Go library to encode and parse [Bencode](https://en.wikipedia.org/wiki/Bencode).

Inspired by the great work done by [@jackpal/bencode-go](https://github.com/jackpal/bencode-go/) and [@marksamman/bencode](https://github.com/marksamman/bencode).

## Examples

```go
import "github.com/stefanovazzocell/bencode"

// Parse a string
bencode.NewParserFromString("10:helloworld").AsString() // "helloworld", nil
// Parse an object
bencode.NewParserFromString("d4:name5:Alice3:agei35ee").AsInterface() // map[string]{}{ "name": "Alice", "age": 35 }, nil

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

- Near-100% code coverage (`make test`, coverage: 97.1%).
- Extensively fuzzed (`make fuzz`, coverage: 99.3%).
- Checked for security issues with [`gosec`](https://github.com/securego/gosec) (`make security`).

### Performance

This library needs to perform well as it might need to encode/decode a large amount of data efficiently. Benchmarks are available with `make bench`.

```
$ make bench
go test -run=^$ -cover -bench .
goos: linux
goarch: amd64
pkg: github.com/stefanovazzocell/bencode
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkEncoder/torrentString-12       12163128                88.54 ns/op
BenchmarkEncoder/complexMap-12           1433679                856.5 ns/op
BenchmarkStringParser/torrentString-12  15259002                75.62 ns/op
BenchmarkStringParser/complexMap-12       910756                 1316 ns/op
```

Furthermore you can profile specific components with the following:
- `make profileEncoder`
- `make profileStringParser`

## Caveats

There are some things to consider
- This library can encode `int` and `uint` as well as all their variations (i.e.: `int64`, `uint16`, ...) but it can only parse numbers of type `int`.
- This library can only encode maps of type `map[string]interface{}` and slices of type `[]interface{}`.
- Additional data after the initial parse will be ignored, unless another parse operation (such as `.AsList()`) is called on the same parser.
