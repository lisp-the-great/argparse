package argparse

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Parser struct {
	Args map[string]*Argument
}

func New(name string) *Parser {
	return &Parser{Args: map[string]*Argument{}}
}

func (p *Parser) Parse(v any, args ...string) {
	vv := reflect.ValueOf(v)
	vt := vv.Type()
	if vt.Kind() != reflect.Ptr || vt.Elem().Kind() != reflect.Struct {
		panic(ErrNotStructPtr)
	}
	vt = vt.Elem()

	for i := 0; i < vt.NumField(); i++ {
		sf := vt.Field(i)

		tag, ok := sf.Tag.Lookup("arg")
		if !ok {
			continue
		}

		arg, err := parseArgument(tag, vv.Elem().Field(i))
		if err != nil {
			panic(err)
		}
		if err = arg.Validate(); err != nil {
			panic(err)
		}
		p.Args[arg.Name] = arg
	}

	for i := 0; i < len(args); i++ {
		var arg *Argument = nil
		if isValidFlag(args[i]) {
			for _, a := range p.Args {
				if a.Flag == args[i][1] {
					arg = a
					break
				}
			}
		} else if args[i][:2] == "--" {
			arg = p.Args[args[i][2:]]
		}

		if arg == nil {
			panic(UnknownArgumentTag(args[i]))
		}
		if arg.Value.Kind() == reflect.Bool {
			arg.Value.SetBool(!arg.DefaultValue.(bool))
		} else if i+1 == len(args) {
			panic(MissingValue(args[i]))
		} else {
			arg.SetValue(args[i+1])
			i++
		}
	}

	for _, arg := range p.Args {
		if arg.Value.IsZero() && arg.DefaultValue != nil {
			arg.Value.Set(reflect.ValueOf(arg.DefaultValue))
		}
		if arg.Required && arg.Value.IsZero() {
			panic(MissingRequiredArgument(arg.Name))
		}
	}
}

func (p *Parser) String() string {
	frags := []string{}
	for k, v := range p.Args {
		frags = append(frags, fmt.Sprintf("%s => %v", k, v))
	}
	return fmt.Sprintf("Parser{%s}", strings.Join(frags, ", "))
}

func Parse(v any) {
	New(os.Args[0]).Parse(v, os.Args[1:]...)
}
