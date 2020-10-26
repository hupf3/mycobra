package mycobra

import (
	"io"
	"strings"
	"text/template"

	flag "github.com/spf13/pflag"
)

var templateFuncs = template.FuncMap{
	"trim": strings.TrimSpace,
}

// hasNoOptDefVal 判断没有‘-’的参数是否存在
func hasNoOptDefVal(name string, fs *flag.FlagSet) bool {
	flag := fs.Lookup(name)
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}

// shortHasNoOptDefVal 判断有‘-’的参数是否存在
func shortHasNoOptDefVal(name string, fs *flag.FlagSet) bool {
	if len(name) == 0 {
		return false
	}

	flag := fs.ShorthandLookup(name[:1])
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}

// stripFlags 解析命令参数列表
func stripFlags(args []string, c *Command) []string {
	if len(args) == 0 {
		return args
	}

	commands := []string{}
	flags := c.Flags()
Loop:
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			// "--" terminates the flags
			break Loop
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !hasNoOptDefVal(s[2:], flags):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !shortHasNoOptDefVal(s[1:], flags):
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				break Loop
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		}
	}

	return commands
}

// tmpl 对数据执行给定的模板文本，将结果写入w
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("usage")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
