package argparse

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotStructPtr = errors.New("not struct pointer")
)

func InvalidTag(s string) error {
	return fmt.Errorf("invalid tag \"%s\"", s)
}

func InvalidType(t reflect.Type) error {
	return fmt.Errorf("invalid type \"%v\"", t)
}

func InvalidFlag(c byte) error {
	return fmt.Errorf("bad flag '%c' must be of format \"-?\" where \"?\" is an alphabet", c)
}

func InvalidName(s string) error {
	return fmt.Errorf("invalid name \"%s\"", s)
}

func MissingRequiredArgument(name string) error {
	return fmt.Errorf("missing required argument \"%s\"", name)
}

func MissingValue(t string) error {
	return fmt.Errorf("missing value for \"%s\"", t)
}

func UnknownArgumentTag(t string) error {
	return fmt.Errorf("unknown argument tag \"%s\"", t)
}
