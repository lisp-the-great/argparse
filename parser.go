package argparse

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type Parser struct {
	Name string
	Args map[string]*Argument
}

func New(name string) *Parser {
	if name == "" {
		panic(ErrEmptyParserName)
	}
	return &Parser{Name: name, Args: map[string]*Argument{}}
}

func (p *Parser) Parse(v any, args ...string) error {
	vv := reflect.ValueOf(v)
	vt := vv.Type()
	if vt.Kind() != reflect.Ptr || vt.Elem().Kind() != reflect.Struct {
		return ErrNotStructPtr
	}
	vt = vt.Elem()

	// Parse target `v` and prepare.
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

	// Seek `-h` or `--help` then show help message and exit
	for _, s := range args {
		if s == "-h" || s == "--help" {
			fmt.Println(p.HelpMessage())
			os.Exit(0)
		}
	}

	// Match `args` with `v`'s fields.
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
		} else {
			return UnknownArgumentTag(args[i])
		}

		if arg == nil {
			panic(UnknownArgumentTag(args[i]))
		}
		if arg.NumValues() == 0 {
			arg.Value.SetBool(!arg.DefaultValue.(bool))
		} else if i+1 == len(args) {
			panic(MissingValue(args[i]))
		} else {
			arg.SetValue(args[i+1])
			i++
		}
	}

	// Validate all required `Argument`s.
	for _, arg := range p.Args {
		if arg.Value.IsZero() && arg.DefaultValue != nil {
			arg.Value.Set(reflect.ValueOf(arg.DefaultValue))
		}
		if arg.Required && arg.Value.IsZero() {
			return MissingRequiredArgument(arg.Name)
		}
	}
	return nil
}

func (p *Parser) String() string {
	frags := []string{}
	for k, v := range p.Args {
		frags = append(frags, fmt.Sprintf("%s => %v", k, v))
	}
	return fmt.Sprintf("Parser{ %s }", strings.Join(frags, ", "))
}

func Parse(v any) {
	if err := New(os.Args[0]).Parse(v, os.Args[1:]...); err != nil {
		panic(err)
	}
}

func (p *Parser) HelpMessage() string {
	args := make([]*Argument, 0, len(p.Args))
	for _, a := range p.Args {
		args = append(args, a)
	}
	sort.Slice(args, func(i, j int) bool {
		return args[i].Name < args[j].Name
	})

	shorts := make([]string, len(args)+1)
	shorts[0] = "[-h]"
	longParts := make([][2]string, len(args)+1)
	longParts[0] = [2]string{"-h,--help", "show this message and exit"}
	for i, a := range args {
		var sh string
		if a.HasFlag() {
			sh = a.Short()
		} else {
			sh = a.Long()
		}
		if a.NumValues() > 0 {
			sh += " " + a.NameUpperCase()
		}
		if !a.Required {
			sh = "[" + sh + "]"
		}
		shorts[i+1] = sh
		pre, suf := a.helpMessage()
		longParts[i+1] = [2]string{pre, suf}
	}

	maxLen := 0
	for _, x := range longParts {
		if l := len(x[0]); l > maxLen {
			maxLen = l
		}
	}
	longs := make([]string, len(longParts))
	for i, x := range longParts {
		longs[i] = x[0] + strings.Repeat(" ", maxLen-len(x[0])) + "\t" + x[1]
	}

	return fmt.Sprintf("Usage: %s %s\n\n\t%s",
		p.Name,
		strings.Join(shorts, " "),
		strings.Join(longs, "\n\t"))
}
