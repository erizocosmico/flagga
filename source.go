package flagga

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Source provides values for the flags.
type Source interface {
	// Open allows the Source to perform some initialization before parsing.
	Open() error
	// Get will try to get the given key from the source and fill the Value.
	Get(key string, dst Value) (bool, error)
	// Close is called after parsing to release all used resources by the
	// source. Close should be tolerant to multiple calls, even if it has
	// not been opened.
	Close() error
}

type envSource string

// EnvPrefix will provide as values the environment variables that match
// the given prefix.
func EnvPrefix(prefix string) Source {
	return envSource(prefix)
}

func (envSource) Open() error  { return nil }
func (envSource) Close() error { return nil }
func (e envSource) Get(key string, dst Value) (bool, error) {
	v, ok := os.LookupEnv(string(e) + key)
	if !ok {
		return false, nil
	}

	if err := dst.Set(v); err != nil {
		return false, err
	}

	return true, nil
}

// FileSource is a Source that reads a file and parses it using a parser
// function.
type FileSource struct {
	File   string
	Parser ParseFunc
	Value  map[string]interface{}
}

// ParseFunc is a function that will parse the given data and put the
// result into the given destination.
type ParseFunc func(data []byte, dst interface{}) error

// NewFileSource returns a Source that will read the given file and use the
// given parser to extract the contents of it.
func NewFileSource(file string, parser ParseFunc) Source {
	return &FileSource{file, parser, nil}
}

type jsonSource struct {
	Source
}

// JSONVia returns a Source that will use a JSON file as a provider of
// flag values.
func JSONVia(file string) Source {
	return &jsonSource{NewFileSource(file, json.Unmarshal)}
}

// Open implements the Source interface.
func (s *FileSource) Open() error {
	var err error
	content, err := ioutil.ReadFile(s.File)
	if err != nil {
		return err
	}

	return s.Parser(content, &s.Value)
}

// Close implements the Source interface.
func (s *FileSource) Close() error {
	return nil
}

// Get implements the Source interface.
func (s *FileSource) Get(key string, dst Value) (bool, error) {
	val, ok := s.Value[key]
	if !ok {
		return false, nil
	}

	if err := dst.Set(val); err != nil {
		return false, err
	}

	return true, nil
}
