package flagga

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Flag is a single flag in the program.
type Flag struct {
	Name       string
	Usage      string
	Value      Value
	Default    interface{}
	Extractors []Extractor
}

// FlagSet is a collection of unique flags.
type FlagSet struct {
	name     string
	parsed   bool
	args     []string
	nonFlags []string
	sources  []Source
	flags    map[string]*Flag
	found    map[string]*Flag
	output   io.Writer
}

// Parse fills the flags with values from the given arguments and sources.
func (fs *FlagSet) Parse(args []string, sources ...Source) error {
	if fs.parsed {
		return nil
	}
	fs.parsed = true
	fs.args = args
	fs.sources = sources

	fs.found = make(map[string]*Flag)
	for {
		var err error
		args, err = fs.parseNext(args)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			break
		}
	}

	defer func() {
		for _, s := range sources {
			_ = s.Close()
		}
	}()

	for _, s := range sources {
		if err := s.Open(); err != nil {
			return err
		}
	}

	// now find the ones that are not filled using other sources
	for name, f := range fs.flags {
		if _, ok := fs.found[name]; !ok {
			var found bool
			for _, e := range f.Extractors {
				if ok, err := e.Get(sources, f.Value); err != nil {
					return err
				} else if ok {
					found = true
					break
				}
			}

			// if no value could be found, just use the default value
			if !found {
				fs.found[name] = f
				if err := f.Value.Set(f.Default); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (fs *FlagSet) parseNext(args []string) ([]string, error) {
	for {
		if len(args) == 0 {
			return nil, nil
		}

		arg, args := args[0], args[1:]
		if len(arg) < 2 || arg[0] != '-' {
			fs.args = append(fs.args, arg)
			// this was not a flag, skip it
			continue
		}

		var name string
		if arg == "--" {
			// -- terminates flags
			fs.args = append(fs.args, args...)
			return nil, nil
		} else if strings.HasPrefix(arg, "--") {
			name = arg[2:]
		} else {
			name = arg[1:]
		}

		if len(name) == 0 || name[0] == '-' || name[0] == '=' {
			return nil, fmt.Errorf("invalid flag syntax: %s", arg)
		}

		idx := strings.IndexRune(name, '=')
		if idx > 0 {
			// has a value
			name, value := name[:idx], name[idx:]
			if len(value) == 0 {
				return nil, fmt.Errorf("invalid flag syntax: %s", arg)
			}

			if err := fs.setValue(name, value); err != nil {
				return nil, err
			}
		} else {
			if f, ok := fs.flags[name]; ok && isBool(f.Value) {
				if err := f.Value.Set(true); err != nil {
					return nil, err
				}
			} else if !ok {
				return nil, fmt.Errorf("unknown flag %s", name)
			}

			if len(args) == 0 {
				return nil, fmt.Errorf("expecting value for flag: %s", name)
			}

			arg, args = args[0], args[1:]
			if strings.HasPrefix(arg, "-") {
				return nil, fmt.Errorf("expecting value for flag: %s", name)
			}

			if err := fs.setValue(name, arg); err != nil {
				return nil, err
			}
		}

		return args, nil
	}
}

func (fs *FlagSet) setValue(name, value string) error {
	f, alreadyFound := fs.found[name]
	if alreadyFound && !isSlice(f.Value) {
		// ignore, we already have a value for this flag
		return nil
	}

	if alreadyFound {
		if err := f.Value.Set(value); err != nil {
			return err
		}
	} else {
		f, ok := fs.flags[name]
		if !ok {
			return fmt.Errorf("unknown flag %s", name)
		}

		fs.found[name] = f
		if err := f.Value.Set(value); err != nil {
			return err
		}
	}

	return nil
}

// Parsed returns whether the flag set has already been parsed.
func (fs *FlagSet) Parsed() bool {
	return fs.parsed
}

func (fs *FlagSet) addFlag(
	name string,
	defaultValue interface{},
	usage string,
	value Value,
	extractors []Extractor,
) {
	if fs.flags == nil {
		fs.flags = make(map[string]*Flag)
	}

	if _, ok := fs.flags[name]; ok {
		panic(fmt.Errorf("flag %s was already defined", name))
	}

	fs.flags[name] = &Flag{
		Name:       name,
		Usage:      usage,
		Default:    defaultValue,
		Value:      value,
		Extractors: extractors,
	}
}

// String adds a new string flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) String(
	name, defaultValue, usage string,
	extractors ...Extractor,
) *string {
	v := new(string)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Int adds a new int flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Int(
	name string,
	defaultValue int,
	usage string,
	extractors ...Extractor,
) *int {
	v := new(int)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Bool adds a new bool flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Bool(
	name string,
	defaultValue bool,
	usage string,
	extractors ...Extractor,
) *bool {
	v := new(bool)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Int64 adds a new int64 flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Int64(
	name string,
	defaultValue int64,
	usage string,
	extractors ...Extractor,
) *int64 {
	v := new(int64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Float adds a new float64 flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Float(
	name string,
	defaultValue float64,
	usage string,
	extractors ...Extractor,
) *float64 {
	v := new(float64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Uint adds a new uint flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Uint(
	name string,
	defaultValue uint,
	usage string,
	extractors ...Extractor,
) *uint {
	v := new(uint)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Uint64 adds a new uint64 flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Uint64(
	name string,
	defaultValue uint64,
	usage string,
	extractors ...Extractor,
) *uint64 {
	v := new(uint64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Duration adds a new time.Duration flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) Duration(
	name string,
	defaultValue time.Duration,
	usage string,
	extractors ...Extractor,
) *time.Duration {
	v := new(time.Duration)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// StringList adds a new []string flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) StringList(
	name string,
	defaultValue []string,
	usage string,
	extractors ...Extractor,
) *[]string {
	v := new([]string)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// IntList adds a new []int flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) IntList(
	name string,
	defaultValue []int,
	usage string,
	extractors ...Extractor,
) *[]int {
	v := new([]int)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Int64List adds a new []int64 flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) Int64List(
	name string,
	defaultValue []int64,
	usage string,
	extractors ...Extractor,
) *[]int64 {
	v := new([]int64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// FloatList adds a new []float64 flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) FloatList(
	name string,
	defaultValue []float64,
	usage string,
	extractors ...Extractor,
) *[]float64 {
	v := new([]float64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// UintList adds a new []uint flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) UintList(
	name string,
	defaultValue []uint,
	usage string,
	extractors ...Extractor,
) *[]uint {
	v := new([]uint)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// Uint64List adds a new []uint64 flag and returns a pointer to the value
// that will be filled once the flag set is parsed.
func (fs *FlagSet) Uint64List(
	name string,
	defaultValue []uint64,
	usage string,
	extractors ...Extractor,
) *[]uint64 {
	v := new([]uint64)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}

// DurationList adds a new []time.Duration flag and returns a pointer to the
// value that will be filled once the flag set is parsed.
func (fs *FlagSet) DurationList(
	name string,
	defaultValue []time.Duration,
	usage string,
	extractors ...Extractor,
) *[]time.Duration {
	v := new([]time.Duration)
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
	return v
}
