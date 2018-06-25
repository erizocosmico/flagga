package flagga

import (
	"fmt"
	"io"
	"os"
	"reflect"
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
	name          string
	description   string
	parsed        bool
	args          []string
	nonFlags      []string
	sources       []Source
	flagOrder     []string
	flags         map[string]*Flag
	found         map[string]*Flag
	out           io.Writer
	errorHandling ErrorHandling

	// Usage prints the usage instructions of the flag set.
	Usage func()
}

// ErrorHandling defines what happens if an error is encountered while parsing
// a flag set.
type ErrorHandling byte

const (
	// ContinueOnError will not halt the program if an error is encountered.
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) after encountering an error.
	ExitOnError
	// PanicOnError will panic after encountering an error.
	PanicOnError
)

// NewFlagSet creates a new flag set with the given name, description and error
// handling policy.
func NewFlagSet(name, description string, errorHandling ErrorHandling) *FlagSet {
	var fs FlagSet
	fs.Init(name, description, errorHandling)
	return &fs
}

// Init initializes the flag set with the given name, description and error
// handling policy.
func (fs *FlagSet) Init(name, description string, errorHandling ErrorHandling) {
	fs.args = nil
	fs.found = make(map[string]*Flag)
	fs.flags = make(map[string]*Flag)
	fs.description = description
	fs.name = name
	fs.errorHandling = errorHandling
}

var exit = os.Exit

// Parse fills the flags with values from the given arguments and sources.
func (fs *FlagSet) Parse(args []string, sources ...Source) error {
	if fs.parsed {
		return nil
	}
	fs.parsed = true
	fs.sources = sources

	if fs.found == nil {
		fs.found = make(map[string]*Flag)
	}
	for {
		var err error
		args, err = fs.parseNext(args)
		if err != nil {
			if err == ErrHelp {
				fs.printUsage()
			} else {
				fs.printError(err)
			}

			switch fs.errorHandling {
			case ContinueOnError:
				return err
			case PanicOnError:
				panic(err)
			case ExitOnError:
				exit(2)
				return nil
			}
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

func (fs *FlagSet) printUsage() {
	if fs.Usage == nil {
		fs.usage()
	} else {
		fs.Usage()
	}
}

func (fs *FlagSet) printError(err error) {
	fmt.Fprint(fs.Output(), err)
	fs.printUsage()
}

// ErrHelp is returned when -h, --h, --help or -help are found.
var ErrHelp = fmt.Errorf("flagga: help requested")

func (fs *FlagSet) usage() {
	if fs.name == "" {
		fmt.Fprint(fs.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.name)
	}

	if fs.description != "" {
		fmt.Fprintf(
			fs.out, "\n  %s\n",
			strings.Replace(fs.description, "\n", "\n  ", -1),
		)
	}

	fmt.Fprint(fs.Output(), "\n")
	fs.PrintDefaults()
}

// PrintDefaults prints all flags with their description and default value.
func (fs *FlagSet) PrintDefaults() {
	for _, name := range fs.flagOrder {
		f := fs.flags[name]
		typ := strings.Replace(
			reflect.TypeOf(f.Default).String(),
			"[]", "list of ", 1,
		)
		fmt.Fprintf(fs.Output(), "  -%s %s\n", name, typ)

		fmt.Fprint(fs.Output(), "  \t")
		if f.Usage != "" {
			fmt.Fprint(
				fs.out,
				strings.Replace(f.Usage, "\n", "\n  \t", -1),
			)
		}

		s, ok := f.Default.(string)
		if !ok || s != "" {
			fmt.Fprintf(fs.Output(), " (default value: %s)\n", prettyValue(f.Default))
		} else {
			fmt.Fprint(fs.Output(), "\n")
		}
	}
}

func prettyValue(v interface{}) string {
	switch v := v.(type) {
	case []string:
		return fmt.Sprintf("[%s]", strings.Join(v, ", "))
	case []int:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case []uint:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case []int64:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case []uint64:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case []float64:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case []time.Duration:
		var parts = make([]string, len(v))
		for i, val := range v {
			parts[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	default:
		return fmt.Sprint(v)
	}
}

func (fs *FlagSet) parseNext(args []string) ([]string, error) {
	for {
		if len(args) == 0 {
			return nil, nil
		}

		arg := args[0]
		args = args[1:]
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

		if name == "h" || name == "help" {
			return nil, ErrHelp
		}

		idx := strings.IndexRune(name, '=')
		if idx > 0 {
			// has a value
			name, value := name[:idx], name[idx+1:]
			if len(value) == 0 {
				return nil, fmt.Errorf("invalid flag syntax: %s", arg)
			}

			if err := fs.setValue(name, value); err != nil {
				return nil, err
			}
		} else {
			f, ok := fs.flags[name]
			if ok && isBool(f.Value) {
				fs.found[name] = f
				if err := f.Value.Set(true); err != nil {
					return nil, err
				}

				return args, nil
			}

			if !ok {
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

		if fs.found == nil {
			fs.found = make(map[string]*Flag)
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

// NFlags returns the number of flags that have been filled.
func (fs *FlagSet) NFlags() int { return len(fs.found) }

// NArg returns the number of arguments that have been found.
func (fs *FlagSet) NArg() int { return len(fs.args) }

// Args returns the arguments that have been found.
func (fs *FlagSet) Args() []string { return fs.args }

// Arg returns the nth argument that has been found.
func (fs *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(fs.args) {
		return ""
	}
	return fs.args[i]
}

// Lookup returns the defined flag with the given name. It will return nil if
// it's not found.
func (fs *FlagSet) Lookup(name string) *Flag { return fs.flags[name] }

// Name returns the given name to this flag set.
func (fs *FlagSet) Name() string { return fs.name }

// Description returns the given description to this flag set.
func (fs *FlagSet) Description() string { return fs.description }

// SetOutput sets the destination writer for the usage and error messages.
func (fs *FlagSet) SetOutput(w io.Writer) { fs.out = w }

// Output returns the destination writer for the usage and error messages. If
// no output was set, the default is os.Stderr.
func (fs *FlagSet) Output() io.Writer {
	if fs.out == nil {
		return os.Stderr
	}
	return fs.out
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

	fs.flagOrder = append(fs.flagOrder, name)

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
	fs.StringVar(v, name, defaultValue, usage, extractors...)
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
	fs.IntVar(v, name, defaultValue, usage, extractors...)
	return v
}

// Bool adds a new bool flag and returns a pointer to the value that will
// be filled once the flag set is parsed.
func (fs *FlagSet) Bool(
	name string,
	usage string,
	extractors ...Extractor,
) *bool {
	v := new(bool)
	fs.BoolVar(v, name, usage, extractors...)
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
	fs.Int64Var(v, name, defaultValue, usage, extractors...)
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
	fs.FloatVar(v, name, defaultValue, usage, extractors...)
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
	fs.UintVar(v, name, defaultValue, usage, extractors...)
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
	fs.Uint64Var(v, name, defaultValue, usage, extractors...)
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
	fs.DurationVar(v, name, defaultValue, usage, extractors...)
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
	fs.StringListVar(v, name, defaultValue, usage, extractors...)
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
	fs.IntListVar(v, name, defaultValue, usage, extractors...)
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
	fs.Int64ListVar(v, name, defaultValue, usage, extractors...)
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
	fs.FloatListVar(v, name, defaultValue, usage, extractors...)
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
	fs.UintListVar(v, name, defaultValue, usage, extractors...)
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
	fs.Uint64ListVar(v, name, defaultValue, usage, extractors...)
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
	fs.DurationListVar(v, name, defaultValue, usage, extractors...)
	return v
}

// StringVar adds a new string flag. When the flag set is parsed it will fill
// the given pointer.
func (fs *FlagSet) StringVar(
	v *string,
	name string,
	defaultValue string,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// IntVar adds a new int flag. When the flag set is parsed it will fill the
// given pointer.
func (fs *FlagSet) IntVar(
	v *int,
	name string,
	defaultValue int,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// UintVar adds a new uint flag. When the flag set is parsed it will fill the
// given pointer.
func (fs *FlagSet) UintVar(
	v *uint,
	name string,
	defaultValue uint,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// Int64Var adds a new int64 flag. When the flag set is parsed it will fill the
// given pointer.
func (fs *FlagSet) Int64Var(
	v *int64,
	name string,
	defaultValue int64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// Uint64Var adds a new uint64 flag. When the flag set is parsedit will fill
// the given pointer.
func (fs *FlagSet) Uint64Var(
	v *uint64,
	name string,
	defaultValue uint64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// BoolVar adds a new bool flag. When the flag set is parsed it will fill the
// given pointer.
func (fs *FlagSet) BoolVar(
	v *bool,
	name string,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, false, usage, NewValue(v), extractors)
}

// FloatVar adds a new float64 flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) FloatVar(
	v *float64,
	name string,
	defaultValue float64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// DurationVar adds a new time.Duration flag. When the flag set is parsed it
// will fill the given pointer.
func (fs *FlagSet) DurationVar(
	v *time.Duration,
	name string,
	defaultValue time.Duration,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// StringListVar adds a new []string flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) StringListVar(
	v *[]string,
	name string,
	defaultValue []string,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// IntListVar adds a new []int flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) IntListVar(
	v *[]int,
	name string,
	defaultValue []int,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// UintListVar adds a new []uint flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) UintListVar(
	v *[]uint,
	name string,
	defaultValue []uint,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// Int64ListVar adds a new []int64 flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) Int64ListVar(
	v *[]int64,
	name string,
	defaultValue []int64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// Uint64ListVar adds a new []uint64 flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) Uint64ListVar(
	v *[]uint64,
	name string,
	defaultValue []uint64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// FloatListVar adds a new []float64 flag. When the flag set is parsed it will
// fill the given pointer.
func (fs *FlagSet) FloatListVar(
	v *[]float64,
	name string,
	defaultValue []float64,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}

// DurationListVar adds a new []time.Duration flag. When the flag set is parsed
// it will fill the given pointer.
func (fs *FlagSet) DurationListVar(
	v *[]time.Duration,
	name string,
	defaultValue []time.Duration,
	usage string,
	extractors ...Extractor,
) {
	fs.addFlag(name, defaultValue, usage, NewValue(v), extractors)
}
