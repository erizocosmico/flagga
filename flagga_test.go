package flagga

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type flag struct {
	name  string
	value interface{}
}

func TestParseNextBool(t *testing.T) {
	var fs FlagSet

	fs.found = make(map[string]*Flag)
	b := fs.Bool("b", "bool value")
	s := fs.String("s", "", "string value")

	remaining, err := fs.parseNext([]string{
		"foo",
		"-b",
		"-s",
		"foo",
	})

	expect(t, fs.args, []string{"foo"})
	expect(t, err, nil)
	expect(t, remaining, []string{"-s", "foo"})
	expect(t, *b, true)
	expect(t, *s, "")
}

func TestParseNextUnknownFlag(t *testing.T) {
	var fs FlagSet

	_, err := fs.parseNext([]string{"-x", "5"})
	expect(t, err, fmt.Errorf("unknown flag x"))
}

func TestParseNextExpectingValue(t *testing.T) {
	var fs FlagSet
	x := fs.Int("x", 0, "int value")

	_, err := fs.parseNext([]string{"-x"})
	expect(t, err, fmt.Errorf("expecting value for flag: x"))
	expect(t, *x, 0)

	_, err = fs.parseNext([]string{"-x", "-n"})
	expect(t, err, fmt.Errorf("expecting value for flag: x"))
	expect(t, *x, 0)
}

func TestParseNextTerminateFlags(t *testing.T) {
	var fs FlagSet

	x := fs.Int("x", 0, "")
	remaining, err := fs.parseNext([]string{"--", "-x", "5"})

	expect(t, err, nil)
	expect(t, remaining, ([]string)(nil))
	expect(t, fs.args, []string{"-x", "5"})
	expect(t, *x, 0)
}

func TestParseNextDoubleDash(t *testing.T) {
	var fs FlagSet

	x := fs.Int("x", 0, "")
	remaining, err := fs.parseNext([]string{"--x", "5"})

	expect(t, err, nil)
	expect(t, remaining, []string{})
	expect(t, *x, 5)
}

func TestParseNextSingleDash(t *testing.T) {
	var fs FlagSet

	x := fs.Int("x", 0, "")
	remaining, err := fs.parseNext([]string{"-x", "5"})

	expect(t, err, nil)
	expect(t, remaining, []string{})
	expect(t, *x, 5)
}

func TestParseNextInlineValue(t *testing.T) {
	var fs FlagSet

	x := fs.Int("x", 0, "")
	remaining, err := fs.parseNext([]string{"-x=5"})

	expect(t, err, nil)
	expect(t, remaining, []string{})
	expect(t, *x, 5)
}

func TestParseNextInvalidFlagSyntax(t *testing.T) {
	var fs FlagSet

	_, err := fs.parseNext([]string{"---x=5"})
	expect(t, err, fmt.Errorf("invalid flag syntax: ---x=5"))

	_, err = fs.parseNext([]string{"-x="})
	expect(t, err, fmt.Errorf("invalid flag syntax: -x="))
}

func TestParse(t *testing.T) {
	var fs FlagSet

	os.Setenv("TEST_B", "env_bar")
	os.Setenv("TEST_C", "env_baz")

	a := fs.String("a", "", "")
	b := fs.String("b", "", "", Env("TEST_B"))
	c := fs.String("c", "", "", Env("TEST_C"))
	d := fs.String("d", "default", "", Env("TEST_D"))

	expect(t, fs.Parsed(), false)

	err := fs.Parse([]string{"-a=foo", "-b=bar", "foo", "bar"}, EnvPrefix(""))
	expect(t, err, nil)

	expect(t, *a, "foo")
	expect(t, *b, "bar")
	expect(t, *c, "env_baz")
	expect(t, *d, "default")
	expect(t, fs.Parsed(), true)
	expect(t, fs.NArg(), 2)
	expect(t, fs.Args(), []string{"foo", "bar"})
	expect(t, fs.Arg(0), "foo")
	expect(t, fs.NFlags(), 3)
}

func TestString(t *testing.T) {
	var fs FlagSet
	x := fs.String("x", "", "")
	expect(t, fs.Parse([]string{"-x=foo"}), nil)
	expect(t, *x, "foo")
}

func TestStringList(t *testing.T) {
	var fs FlagSet
	x := fs.StringList("x", nil, "")
	expect(t, fs.Parse([]string{"-x=foo", "-x=bar", "-x=baz"}), nil)
	expect(t, *x, []string{"foo", "bar", "baz"})
}

func TestBool(t *testing.T) {
	var fs FlagSet
	x := fs.Bool("x", "")
	expect(t, fs.Parse([]string{"-x"}), nil)
	expect(t, *x, true)
}

func TestInt(t *testing.T) {
	var fs FlagSet
	x := fs.Int("x", 0, "")
	expect(t, fs.Parse([]string{"-x=5"}), nil)
	expect(t, *x, 5)
}

func TestIntList(t *testing.T) {
	var fs FlagSet
	x := fs.IntList("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1", "-x=2", "-x=3"}), nil)
	expect(t, *x, []int{1, 2, 3})
}

func TestUint(t *testing.T) {
	var fs FlagSet
	x := fs.Uint("x", 0, "")
	expect(t, fs.Parse([]string{"-x=5"}), nil)
	expect(t, *x, uint(5))
}

func TestUintList(t *testing.T) {
	var fs FlagSet
	x := fs.UintList("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1", "-x=2", "-x=3"}), nil)
	expect(t, *x, []uint{1, 2, 3})
}

func TestInt64(t *testing.T) {
	var fs FlagSet
	x := fs.Int64("x", 0, "")
	expect(t, fs.Parse([]string{"-x=5"}), nil)
	expect(t, *x, int64(5))
}

func TestInt64List(t *testing.T) {
	var fs FlagSet
	x := fs.Int64List("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1", "-x=2", "-x=3"}), nil)
	expect(t, *x, []int64{1, 2, 3})
}

func TestUint64(t *testing.T) {
	var fs FlagSet
	x := fs.Uint64("x", 0, "")
	expect(t, fs.Parse([]string{"-x=5"}), nil)
	expect(t, *x, uint64(5))
}

func TestUint64List(t *testing.T) {
	var fs FlagSet
	x := fs.Uint64List("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1", "-x=2", "-x=3"}), nil)
	expect(t, *x, []uint64{1, 2, 3})
}

func TestDuration(t *testing.T) {
	var fs FlagSet
	x := fs.Duration("x", 0, "")
	expect(t, fs.Parse([]string{"-x=5s"}), nil)
	expect(t, *x, 5*time.Second)
}

func TestDurationList(t *testing.T) {
	var fs FlagSet
	x := fs.DurationList("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1s", "-x=2s", "-x=3s"}), nil)
	expect(t, *x, []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second})
}

func TestFloat(t *testing.T) {
	var fs FlagSet
	x := fs.Float("x", 0, "")
	expect(t, fs.Parse([]string{"-x=3.14"}), nil)
	expect(t, *x, 3.14)
}

func TestFloatList(t *testing.T) {
	var fs FlagSet
	x := fs.FloatList("x", nil, "")
	expect(t, fs.Parse([]string{"-x=1.1", "-x=2.2", "-x=3.3"}), nil)
	expect(t, *x, []float64{1.1, 2.2, 3.3})
}

func TestUsage(t *testing.T) {
	fs := NewFlagSet("foo", "first line\nsecond line", ContinueOnError)
	fs.Bool("a", "flag a")
	fs.String("b", "", "flag b")
	fs.IntList("c", []int{1, 2, 3}, "flag c\nis multiline")

	var buf bytes.Buffer
	fs.SetOutput(&buf)

	fs.printUsage()

	expected := "Usage of foo:\n\n" +
		"  first line\n" +
		"  second line\n" +
		"\n" +
		"  -a bool\n" +
		"  \tflag a (default value: false)\n" +
		"  -b string\n" +
		"  \tflag b\n" +
		"  -c list of int\n" +
		"  \tflag c\n" +
		"  \tis multiline (default value: [1, 2, 3])\n"

	expect(t, buf.String(), expected)

	buf.Reset()
	fs.Usage = func() {
		buf.WriteString("hello")
	}

	fs.printUsage()
	expect(t, buf.String(), "hello")
}

func TestErrorHandling(t *testing.T) {
	t.Run("ContinueOnError", func(t *testing.T) {
		var buf bytes.Buffer
		var fs = NewFlagSet("", "", ContinueOnError)
		fs.SetOutput(&buf)

		err := fs.Parse([]string{"-x"})
		expect(t, err, fmt.Errorf("unknown flag x"))
		expect(t, len(buf.Bytes()) > 0, true)
	})

	t.Run("PanicOnError", func(t *testing.T) {
		var buf bytes.Buffer

		defer func() {
			if r := recover(); r != nil {
				expect(t, len(buf.Bytes()) > 0, true)
				expect(t, r, fmt.Errorf("unknown flag x"))
			} else {
				t.Error("expecting a panic")
			}
		}()

		var fs = NewFlagSet("", "", PanicOnError)
		fs.SetOutput(&buf)

		fs.Parse([]string{"-x"})
	})

	t.Run("ExitOnError", func(t *testing.T) {
		var buf bytes.Buffer
		var code int
		exit = func(i int) {
			code = i
		}

		var fs = NewFlagSet("", "", ExitOnError)
		fs.SetOutput(&buf)

		fs.Parse([]string{"-x"})
		expect(t, len(buf.Bytes()) > 0, true)
		expect(t, code, 2)

		exit = os.Exit
	})
}

func TestHelp(t *testing.T) {
	testCases := []string{"-h", "-help", "--help", "--h"}

	for _, tt := range testCases {
		t.Run(tt, func(t *testing.T) {
			expect(t, new(FlagSet).Parse([]string{tt}), ErrHelp)
		})
	}
}

func TestPrettyValue(t *testing.T) {
	testCases := []struct {
		val      interface{}
		expected string
	}{
		{"foo", "foo"},
		{int(1), "1"},
		{int64(1), "1"},
		{uint(1), "1"},
		{uint64(1), "1"},
		{float64(3.14), "3.14"},
		{true, "true"},
		{1 * time.Second, "1s"},
		{[]string{"a", "b", "c"}, "[a, b, c]"},
		{[]int{1, 2, 3}, "[1, 2, 3]"},
		{[]uint{1, 2, 3}, "[1, 2, 3]"},
		{[]int64{1, 2, 3}, "[1, 2, 3]"},
		{[]uint64{1, 2, 3}, "[1, 2, 3]"},
		{[]float64{1.1, 2.2, 3.3}, "[1.1, 2.2, 3.3]"},
		{[]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}, "[1s, 2s, 3s]"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("(%T)(%v)", tt.val, tt.val), func(t *testing.T) {
			expect(t, prettyValue(tt.val), tt.expected)
		})
	}
}

func expect(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, got: %v", expected, actual)
	}
}
