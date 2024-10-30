package argparse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Argument struct {
	// Flag indicates the short repr of an argument.
	Flag byte

	// Name is the long repr of an argument.
	Name string

	// Value is the actual value to be used.
	Value reflect.Value

	// Required indicates this argument is must passed by user or not.
	Required bool

	// DefaultValue is parsed from `default`.
	DefaultValue any

	// Help is the help message.
	Help string
}

func parseArgument(s string, v reflect.Value) (*Argument, error) {
	arg := &Argument{Value: v}

	i, j := 0, 0
	for j < len(s) {
		if s[j] == ';' {
			if err := parseFrag(s[i:j], arg); err != nil {
				return nil, err
			}
			i = j + 1
		}
		j++
	}
	if err := parseFrag(strings.TrimSpace(s[i:]), arg); err != nil {
		return nil, err
	}

	return arg, arg.Validate()
}

func parseFrag(s string, a *Argument) error {
	switch {
	case isValidFlag(s):
		a.Flag = s[1]
	case s[:2] == "--":
		if name := s[2:]; isValidName(name) {
			a.Name = name
		} else {
			return InvalidName(name)
		}
	case s == "required":
		a.Required = true
	case strings.HasPrefix(s, "default="):
		a.setDefaultValue(s[8:])
	case strings.HasPrefix(s, "help="):
		a.Help = s[5:]
	default:
		UnknownArgumentTag(s)
	}
	return nil
}

func isValidFlag(s string) bool {
	return len(s) == 2 && s[0] == '-' && isValidFlagChar(s[1])
}

func isValidFlagChar(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func (a *Argument) HasFlag() bool {
	return isValidFlagChar(a.Flag)
}

func isValidName(name string) bool {
	return len(name) != 0
}

func (a *Argument) Validate() error {
	if !isValidName(a.Name) {
		return InvalidName(a.Name)
	}
	return nil
}

func (a *Argument) setDefaultValue(v string) (err error) {
	a.DefaultValue, err = parseValue(v, a.Value.Type())
	return
}

func (a *Argument) SetValue(v string) error {
	vv, err := parseValue(v, a.Value.Type())
	if err != nil {
		return err
	}
	a.Value.Set(reflect.ValueOf(vv))
	return nil
}

func parseValue(v string, t reflect.Type) (c any, err error) {
	switch t.Kind() {
	case reflect.Bool:
		c, err = strconv.ParseBool(v)

	case reflect.Int:
		c, err = strconv.Atoi(v)
	case reflect.Int8:
		var x int64
		x, err = strconv.ParseInt(v, 10, 8)
		c = int8(x)
	case reflect.Int16:
		var x int64
		x, err = strconv.ParseInt(v, 10, 16)
		c = int16(x)
	case reflect.Int32:
		var x int64
		x, err = strconv.ParseInt(v, 10, 32)
		c = int32(x)
	case reflect.Int64:
		c, err = parseInt64(t, v)

	case reflect.Uint:
		var x uint64
		x, err = strconv.ParseUint(v, 10, 64)
		c = uint(x)
	case reflect.Uint8:
		var x uint64
		x, err = strconv.ParseUint(v, 10, 8)
		c = uint8(x)
	case reflect.Uint16:
		var x uint64
		x, err = strconv.ParseUint(v, 10, 16)
		c = uint16(x)
	case reflect.Uint32:
		var x uint64
		x, err = strconv.ParseUint(v, 10, 32)
		c = uint32(x)
	case reflect.Uint64:
		var x uint64
		x, err = strconv.ParseUint(v, 10, 64)
		c = uint64(x)

	case reflect.Float32:
		var x float64
		x, err = strconv.ParseFloat(v, 32)
		c = float32(x)
	case reflect.Float64:
		c, err = strconv.ParseFloat(v, 64)

	case reflect.String:
		c = v

	default:
		err = InvalidType(t)
	}
	return
}

func qualTypeName(t reflect.Type) string {
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func parseInt64(t reflect.Type, v string) (any, error) {
	if qualTypeName(t) == "time.Duration" {
		return time.ParseDuration(v)
	}
	return strconv.ParseInt(v, 10, 64)
}

func (a *Argument) String() string {
	name := a.Name
	if a.HasFlag() {
		name = fmt.Sprintf("%s[%c]", a.Name, a.Flag)
	}

	frags := []string{
		fmt.Sprintf("%s = %v", name, a.Value.Interface()),
	}
	if a.Required {
		frags = append(frags, "(*)")
	}
	if a.DefaultValue != nil {
		frags = append(frags, fmt.Sprintf("[%v]", a.DefaultValue))
	}

	return fmt.Sprintf("Argument{%s}", strings.Join(frags, " "))
}

func (a *Argument) Short() string {
	return fmt.Sprintf("-%c", a.Flag)
}

func (a *Argument) Long() string {
	return "--" + a.Name
}

func (a *Argument) NameUpperCase() string {
	return strings.ToUpper(a.Name)
}

func (a *Argument) HelpMessage() string {
	frags := []string{}
	if a.HasFlag() {
		frags = append(frags, a.Short()+",")
	}
	frags = append(frags, a.Long(), " ", a.NameUpperCase(), "\t", a.Help)
	if a.DefaultValue != nil {
		frags = append(frags, fmt.Sprintf("\t[default: %v]", a.DefaultValue))
	}
	return strings.Join(frags, "")
}
