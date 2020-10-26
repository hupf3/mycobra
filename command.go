package mycobra

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// Command 与命令相关的成员变量和函数的结构体
type Command struct {
	Use string // 命令的名称

	Short string // 命令短介绍

	Long string // 长命令介绍

	Example string // 如何使用命令的例子

	usageTemplate string // 命令使用模板

	flags *flag.FlagSet // 全部命令参数

	lflags *flag.FlagSet // 仅针对某个命令的参数，局部的参数

	gflags *flag.FlagSet // 针对所有命令的全局参数

	flagErrorBuf *bytes.Buffer // 包含了 pflag 的错误信息

	commands []*Command // 该命令的子命令列表

	parent *Command // 该命令的父命令

	usageFunc func(*Command) error // 命令的使用介绍

	Run func(cmd *Command, args []string) // 执行命令的函数
}

// getFaGflags 继承了父命令的全局命令
func (c *Command) getFaGflags() {
	// 根命令则无行为
	if c.Parent() == nil {
		return
	}

	c.gflags = c.Parent().getGflags()
}

// ParseFlags 解析持久标志树和局部标志
func (c *Command) ParseFlags(args []string) error {
	if c.flagErrorBuf == nil {
		c.flagErrorBuf = new(bytes.Buffer)
	}
	beforeBufferLen := c.flagErrorBuf.Len()

	c.getFaGflags()
	err := c.Flags().Parse(args)
	if c.flagErrorBuf.Len()-beforeBufferLen > 0 && err == nil {
		fmt.Println(c.flagErrorBuf.String())
	}

	return err
}

// execute 根据指令执行命令
func (c *Command) execute(a []string) error {
	err := c.ParseFlags(a)
	if err != nil {
		return err
	}
	c.Run(c, a)

	return nil
}

// ExecuteC 执行命令并抛出异常
func (c *Command) ExecuteC() (err error) {
	args := os.Args
	cmd, flags, err := c.Find(args)
	if err == errors.New("Help") {
		cmd.Usage()
		return nil
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	return cmd.execute(flags)
}

// Execute 执行命令
func (c *Command) Execute() error {
	err := c.ExecuteC()
	return err
}

// Name 获取命令的名字
func (c *Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// GetLong 返回命令的长介绍
func (c *Command) GetLong() string {
	return c.Long
}

// GetShort 返回命令的短介绍
func (c *Command) GetShort() string {
	return c.Short
}

// Root 返回该命令的根命令
func (c *Command) Root() *Command {
	pa := c
	for pa.parent != nil {
		pa = c.parent
	}
	return pa
}

// Commands 获取该命令的子命令
func (c *Command) Commands() []*Command {
	return c.commands
}

// Parent 返回当前命令的父命令
func (c *Command) Parent() *Command {
	return c.parent
}

// setgflags 设置全局参数
func (c *Command) setGflags(flags *flag.FlagSet) {
	c.gflags = flags
}

// getGflags 获取全局的flags
func (c *Command) getGflags() *flag.FlagSet {
	c.getFaGflags()

	if c.gflags == nil {
		c.gflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.gflags.SetOutput(c.flagErrorBuf)
	}

	return c.gflags
}

// LocalFlags 获取子命令
func (c *Command) LocalFlags() *flag.FlagSet {
	c.getFaGflags()

	if c.lflags == nil {
		c.lflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.lflags.SetOutput(c.flagErrorBuf)
	}

	return c.lflags
}

// Flags 返回完整的命令参数
func (c *Command) Flags() *flag.FlagSet {
	// 获取父命令全局命令
	c.getFaGflags()

	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.flags.SetOutput(c.flagErrorBuf)
	}
	// 添加命令
	c.flags.AddFlagSet(c.lflags)
	c.flags.AddFlagSet(c.gflags)

	return c.flags
}

// AddCommand 添加子命令
func (c *Command) AddCommand(cmds ...*Command) {
	for i, x := range cmds {
		if cmds[i] == c {
			panic("Command can't be a child of itself")
		}
		cmds[i].parent = c
		c.commands = append(c.commands, x)
	}
}

// innerFind 查找要执行的子命令
func innerFind(cmd *Command, innerArgs []string) (*Command, []string, error) {
	if innerArgs[0] != cmd.Name() {
		return cmd, nil, errors.New("The command does not exit")
	}

	argsWOflags := stripFlags(innerArgs[1:], cmd)

	if len(argsWOflags) > 0 && argsWOflags[0] == "help" {
		return cmd, nil, errors.New("Help")
	}

	if len(argsWOflags) == 0 {
		return cmd, innerArgs[1:], nil
	}

	sub := argsWOflags[0]

	subCmd := cmd.findSubcmd(sub)
	if subCmd == nil {
		return cmd, nil, errors.New("The command does not exit")
	}

	return innerFind(subCmd, innerArgs[1:])
}

// Find 找到要执行的子命令
func (c *Command) Find(args []string) (*Command, []string, error) {
	commandFound, flags, err := innerFind(c, args)
	if err == errors.New("Help") {
		return commandFound, []string{}, errors.New("Help")
	}
	if err != nil {
		return commandFound, flags, err
	}
	return commandFound, flags, nil
}

// findSubcmd 寻找特定子命令
func (c *Command) findSubcmd(_name string) *Command {
	for _, cmd := range c.commands {
		if cmd.Name() == _name {
			return cmd
		}
	}

	return nil
}

// Runnable 判断命令本身能否执行
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// IsAvailableCommand 判断是否可用作非帮助命令
func (c *Command) IsAvailableCommand() bool {
	if c.Runnable() || c.HasAvailableSubCommands() {
		return true
	}
	return false
}

// HasAvailableSubCommands 判断该命令是否有有效的子命令
func (c *Command) HasAvailableSubCommands() bool {
	for _, sub := range c.commands {
		if sub.IsAvailableCommand() {
			return true
		}
	}
	return false
}

// HasSubCommands 判断命令是否有子命令
func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

// HasParent 判断命令是否有父命令
func (c *Command) HasParent() bool {
	return c.parent != nil
}

// HasAvailableFlags 判断命令是否包含参数
func (c *Command) HasAvailableFlags() bool {
	c.getFaGflags() // 先获取父亲命令

	return c.Flags().HasAvailableFlags()
}

// HasAvailableGlobalFlags 判断命令是否存在全局参数
func (c *Command) HasAvailableGlobalFlags() bool {
	c.getFaGflags() // 先获取父亲命令

	return c.getGflags().HasAvailableFlags()
}

// HasAvailableLocalFlags 判断命令是否存在局部参数
func (c *Command) HasAvailableLocalFlags() bool {
	return c.LocalFlags().HasAvailableFlags()
}

// UsageFunc 返回为此命令设置的函数
func (c *Command) UsageFunc() (f func(*Command) error) {
	if c.usageFunc != nil {
		return c.usageFunc
	}
	if c.HasParent() {
		return c.Parent().UsageFunc()
	}
	return func(c *Command) error {
		c.getFaGflags()
		err := tmpl(os.Stdout, c.UsageTemplate(), c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return err
	}
}

// Usage 打印命令的使用方法
func (c *Command) Usage() error {
	return c.UsageFunc()(c)
}

// CommandPath 返回到达该子命令的路径
func (c *Command) CommandPath() string {
	if c.HasParent() {
		return c.Parent().CommandPath() + " " + c.Name()
	}
	return c.Name()
}

// UseLine 输出该命令的完整描述
func (c *Command) UseLine() string {
	var useline string
	if c.HasParent() {
		useline = c.parent.CommandPath() + " " + c.Use
	} else {
		useline = c.Use
	}

	if c.HasAvailableFlags() && !strings.Contains(useline, "[flags]") {
		useline += " [flags]"
	}
	return useline
}

// UsageTemplate 返回命令的使用模板
func (c *Command) UsageTemplate() string {
	if c.usageTemplate != "" {
		return c.usageTemplate
	}

	if c.HasParent() {
		return c.parent.UsageTemplate()
	}
	return `
{{.GetLong}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{.Name}}: {{.GetShort}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
LocalFlags:
  {{.LocalFlags.FlagUsages}}
{{end}}{{if .HasAvailableGlobalFlags}}
GlobalFlags:
  {{.GlobalFlags.FlagUsages}}
{{end}} {{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
