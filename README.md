# flagga [![GoDoc](https://godoc.org/github.com/erizocosmico/flagga?status.svg)](https://godoc.org/github.com/erizocosmico/flagga) [![Build Status](https://travis-ci.org/erizocosmico/flagga.svg?branch=master)](https://travis-ci.org/erizocosmico/flagga) [![codecov](https://codecov.io/gh/erizocosmico/flagga/branch/master/graph/badge.svg)](https://codecov.io/gh/erizocosmico/flagga) [![Go Report Card](https://goreportcard.com/badge/github.com/erizocosmico/flagga)](https://goreportcard.com/report/github.com/erizocosmico/flagga) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**flagga** is an extensible Go library for handling program configuration using (but not limited to) command line arguments, environment variables and JSON.

This idea and API come from [Peter Bourgon's *Go for Industrial Programming*](http://peter.bourgon.org/go-for-industrial-programming/#program-configuration) talk at Gophercon Iceland 2018.

## Status

Work in Progress. The library is missing functionality and tests at the moment. The public API may change.

### TODO

- Tests for `FlagSet`.
- Usage text generation.
- Automatic `help` flag detection.
- YAML `Extractor` and `Source`.
- TOML `Extractor` and `Source`.

## Goals

- Be able to configure a program with different sources that have different priorities.
- Be extensible so anyone can extend the API to provide different sources to get their configuration from (yaml, toml, database?, ...).
- Provide as much functionality as the Go standard `flag` package.

## Install

```
go get github.com/erizocosmico/flagga
```

Or use your preferred dependency manager such as dep or vgo.

## Usage


```go
var fs flagga.FlagSet

db := fs.String("db", defaultDBURI, "database connection string", flagga.Env("DBURI"))
users := fs.StringList("users", nil, "list of allowed users", flagga.JSON("users"))

err := fs.Parse(os.Args[1:], flagga.JSONVia("config.json"), flagga.EnvPrefix("MYAPP_"))
if err != nil {
    // handle err
}

fmt.Println(*db) // Outputs: "user@localhost:1234/foo"
fmt.Println(strings.Join(*users, ", ")) // Outputs: "jane, joe, alice"
```

To get the previous results we can invoke the program in the following ways:

```bash
echo '{"users":["jane", "joe", "alice"]}' > config.json
./myprogram -db=user@localhost:1234/foo -users=jane -users=joe -users=alice
MYAPP_DBURI=user@localhost:1234/foo ./myprogram
```

### Priority of sources

CLI flags **always** have priority over environment variables or JSON keys. If a flag is provided using the command line flags, no other sources will be checked for that variable.

The rest of the priorities depend of the order in which the sources are passed to the `Parse` method. For example, `fs.Parse(os.Args, flagga.EnvPrefix("FOO_"), flagga.JSONVia("cfg"))` gives more priority to environment variables than to the JSON configuration.

## Custom `Source`s and `Extractor`s

You can implement your own `Source`s and `Extractor`s in case your configuration is in a different format. Check out the `Source` and `Extractor` interfaces in the package documentation.

### Reference

- [Go `flag` package](http://golang.org/pkg/flag)
- [*Go for Industrial Programming* by Peter Bourgon](http://peter.bourgon.org/go-for-industrial-programming/#program-configuration)

## License

MIT, see [LICENSE](/LICENSE)
