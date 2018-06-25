package flagga

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
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

type fileSource struct {
	file   string
	parser func([]byte, interface{}) error
	value  map[string]interface{}
}

type jsonSource struct {
	fileSource
}

type yamlSource struct {
	fileSource
}

type tomlSource struct {
	fileSource
}

// JSONVia returns a Source that will use a JSON file as a provider of
// flag values.
func JSONVia(file string) Source {
	return &jsonSource{fileSource{file, json.Unmarshal, nil}}
}

// YAMLVia returns a Source that will use a YAML file as a provider of
// flag values.
func YAMLVia(file string) Source {
	return &yamlSource{fileSource{file, yaml.Unmarshal, nil}}
}

// TOMLVia returns a Source that will use a TOML file as a provider of
// flag values.
func TOMLVia(file string) Source {
	return &tomlSource{fileSource{file, toml.Unmarshal, nil}}
}

func (s *fileSource) Open() error {
	var err error
	content, err := ioutil.ReadFile(s.file)
	if err != nil {
		return err
	}

	return s.parser(content, &s.value)
}

func (s *fileSource) Close() error {
	return nil
}

func (s *fileSource) Get(key string, dst Value) (bool, error) {
	val, ok := s.value[key]
	if !ok {
		return false, nil
	}

	if err := dst.Set(val); err != nil {
		return false, err
	}

	return true, nil
}
