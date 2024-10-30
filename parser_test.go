package argparse_test

import (
	"testing"
	"time"

	"github.com/lisp-the-great/argparse"
)

type Arguments struct {
	Name    string        `arg:"-n;--name;required;help=name of the person"`
	MaxAge  int           `arg:"-a;--max-age;default=18;help=and, the max age"`
	Gender  bool          `arg:"--gender;default=true;help=hello boys and girls"`
	Latency time.Duration `arg:"-L;--latency;default=3m59s200ms;help=no actually meaning"`
}

func TestParse(t *testing.T) {
	arg := new(Arguments)
	ap := argparse.New("test")
	ap.Parse(arg, "-a", "24", "--gender", "--name", "Kamala", "-L", "2h3m4s5ms")
	t.Logf("%+v", ap)
	for _, a := range ap.Args {
		t.Logf("%+v", a)
	}
	t.Logf("%+v", arg)
}

func TestHelp(t *testing.T) {
	arg := &Arguments{}
	ap := argparse.New("testHelp")
	ap.Parse(arg, "--help")
}
