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

type jsonSource struct {
	file string
	json map[string]interface{}
}

// JSONVia returns a Source that will use a JSON file as a provider of
// flag values.
func JSONVia(file string) Source {
	return &jsonSource{file, nil}
}

func (s *jsonSource) Open() error {
	var err error
	content, err := ioutil.ReadFile(s.file)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &s.json)
}

func (s *jsonSource) Close() error {
	return nil
}

func (s *jsonSource) Get(key string, dst Value) (bool, error) {
	val, ok := s.json[key]
	if !ok {
		return false, nil
	}

	if err := dst.Set(val); err != nil {
		return false, err
	}

	return true, nil
}
