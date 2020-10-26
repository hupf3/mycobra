# 支持子命令命令行程序支持包开发

## 概述

命令行实用程序并不是都象 cat、more、grep 是简单命令。[go](https://go-zh.org/cmd/go/) 项目管理程序，类似 java 项目管理 [maven](http://maven.apache.org/)、Nodejs 项目管理程序 [npm](https://www.npmjs.com/)、git 命令行客户端、 docker 与 kubernetes 容器管理工具等等都是采用了较复杂的命令行。即一个实用程序同时支持多个子命令，每个子命令有各自独立的参数，命令之间可能存在共享的代码或逻辑，同时随着产品的发展，这些命令可能发生功能变化、添加新命令等。因此，符合 [OCP 原则](https://en.wikipedia.org/wiki/Open/closed_principle) 的设计是至关重要的编程需求。

## 课程任务

- 了解 [Cobra包](https://github.com/spf13/cobra)，使用 cobra 命令行生成一个简单的带子命令的命令行程序
- 模仿 `cobra.Command` 编写一个 myCobra 库
- 将带子命令的命令行处理程序的 `import ("github.com/spf13/cobra")` 改为 `import (corbra "gitee.com/yourId/yourRepo")`
- 使得命令行处理程序修改代价最小，即可正常运行

**任务要求**

1. 核心任务，就是模仿 cobra 库的 command.go 重写一个 Command.go
   - 仅允许使用的第三方库 `flag "github.com/spf13/pflag"`
   - 可以参考、甚至复制原来的代码
   - 必须实现简化版的 `type Command struct` 定义和方法
   - 不一定完全兼容 `github.com/spf13/cobra`
   - 可支持简单带子命令的命令行程序开发
2. 包必须包括以下内容：
   - 生成的中文 api 文档
   - 有较好的 Readme 文件，包括一个简单的使用案例
   - 每个go文件必须有对应的测试文件

## 博客地址

[传送门](https://blog.csdn.net/qq_43267773/article/details/109300701)

## 设计说明

### 获取包

输入以下的命令即可获取我实现的 mycobra 包

`go get github.com/hupf3/mycobra`

或者在 src 的相应目录下输入以下命令

`git clone https://github.com/hupf3/mycobra.git `

`go install`

### 使用包

在代码中直接进行引用：

`import "github.com/hupf3/mycobra"`

如果习惯用 cobra 库的也可以用以下的命令：

`import cobra "github.com/hupf3/mycobra"`

### 简单说明

此次作业 mycobra 包的设计实现，大部分是参考原作者的 cobra 包进行实现的，所以有的部分设计会与原作者重合。主要参考了原作者两个代码文件：[command.go](https://github.com/spf13/cobra/blob/master/command.go)，和 [cobra.go](https://github.com/spf13/cobra/blob/master/cobra.go)

### 包代码文件结构

本次包的实现主要有 5 个代码文件，如下所示：

<img src="https://img-blog.csdnimg.cn/2020102620393350.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom: 50%;" />

- `bench_test.go`：该文件是代码的基准测试，用来测试代码执行的时间

- `cobra.go`：该文件是基于原作者 `cobra.go` 的代码文件进行修改实现的

  - `hasNoOptDefVal()`：实现了判断没有 ‘-’ 的参数是否实现

    ```go
    // hasNoOptDefVal 判断没有‘-’的参数是否存在
    func hasNoOptDefVal(name string, fs *flag.FlagSet) bool {
    	flag := fs.Lookup(name)
    	if flag == nil {
    		return false
    	}
    	return flag.NoOptDefVal != ""
    }
    ```

  - `shortHasNoOptDefVal()`：实现判断有 ‘-’ 的参数是否存在

    ```go
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
    ```

  - `stripFlags()`：实现了解析命令参数列表

    ```go
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
    ```

  - `tmpl()`：实现对数据执行给定的模板文本，将结果写入w

    ```go
    // tmpl 对数据执行给定的模板文本，将结果写入w
    func tmpl(w io.Writer, text string, data interface{}) error {
    	t := template.New("usage")
    	t.Funcs(templateFuncs)
    	template.Must(t.Parse(text))
    	return t.Execute(w, data)
    }
    ```

- `cobra_test.go`：是基于 `cobra.go` 文件的测试文件，主要用于测试该文件的函数实现是否正确

- `command.go`：该代码文件主要是实现了命令和命令相关操作的函数，由于实现的代码量较大，就不在此一一详述，只挑几个重要的结构体和函数进行说明，在 API文档中有每个函数的详细说明。

  <img src="https://img-blog.csdnimg.cn/20201026210325375.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

  - `type Command struct` 与命令相关的成员变量和函数的结构体

    ```go
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
    ```

  - `func (c *Command) AddCommand(cmds ...*Command)`：添加该命令下的子命令

    ```go
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
    ```

  - `func (c *Command) UsageTemplate() string`：返回命令的使用模板(注意：在实现此函数的过程中我发现了一个容易出错的细节，就是在'{{}}'中引用的函数的首字母一定是大写的，如果是小写的会报错)

    ```go
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
    
    Examples:
    	{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
    
    Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
    	{{.Name}}: {{.GetShort}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
    
    Flags:
    	{{.LocalFlags.FlagUsages}}{{end}}{{if .HasAvailableGlobalFlags}}
    
    GlobalFlags:
    	{{.GlobalFlags.FlagUsages}}{{end}} {{if .HasAvailableSubCommands}}
    
    Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
    `
    }
    ```

  - `func (c *Command) Find(args []string) (*Command, []string, error)` 找到要执行的子命令

    ```go
    // innerFind 查找要执行的子命令
    func innerFind(cmd *Command, innerArgs []string) (*Command, []string, error) {
    	if innerArgs[0] != cmd.Name() {
    		return cmd, nil, errors.New("The command does not exit")
    	}
    
    	argsWOflags := stripFlags(innerArgs[1:], cmd)
    
    	if len(argsWOflags) > 0 && argsWOflags[0] == "help" {
    		return cmd, nil, errh
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
    	if err == errh {
    		return commandFound, []string{}, errh
    	}
    	if err != nil {
    		return commandFound, flags, err
    	}
    	return commandFound, flags, nil
    }
    ```

- `command_test.go`：是基于 `command.go` 文件的测试文件，主要用于测试该文件的函数实现是否正确

## 单元测试

### cobra_test.go

此测试文件是基于 `cobra.go` 文件的测试文件，主要用于测试该文件的函数实现是否正确：

```go
package mycobra

import (
	"reflect"
	"testing"
)

func Test_StripFlags(t *testing.T) {
	test := &Command{
		Use:     "test",
		Short:   "test",
		Long:    "test",
		Example: "test",
	}
	// 三种带参数的方法
	test.Flags().StringP("test1", "a", "", "")
	test.Flags().StringP("test2", "b", "", "")
	test.Flags().StringP("test3", "c", "", "")

	input := []string{"-a", "-b1", "-c=1", "test1", "test2", "test3"}
	r := stripFlags(input, test)
	expected := []string{"test1", "test2", "test3"}

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("expected '%q' but got '%q'", expected, r)
	}
}
```

写好代码文件后开始进行测试，测试的结果如下：

<img src="https://img-blog.csdnimg.cn/20201026211150154.png#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

通过上面的结果可以得知函数实现的正确，且通过了测试

### command_test.go

此文件是基于 `command.go` 文件的测试文件，主要用于测试该文件的函数实现是否正确

首先在测试文件中定义两个 Command 变量方便后续的测试：

```go
// 根命令
var test1 = &Command{
	Use:     "test1",
	Short:   "test1",
	Long:    "test1",
	Example: "test1",
}

// 子命令
var test2 = &Command{
	Use:     "test2",
	Short:   "test2",
	Long:    "test2",
	Example: "test2",
}
```

- `func Test_ParseFlags(t *testing.T)`：该函数测试了命令带参数的情况

  ```go
  func Test_ParseFlags(t *testing.T) {
  	// 三种带参数的方法
  	test2.Flags().StringP("test1", "a", "", "")
  	test2.Flags().StringP("test2", "b", "", "")
  	test2.Flags().StringP("test3", "c", "", "")
  
  	args := []string{"-a1", "-b=1"}
  	test2.ParseFlags(args)
  	r1, _ := test2.Flags().GetString("test1")
  	r2, _ := test2.Flags().GetString("test2")
  	e1, e2 := "1", "1"
  
  	if r1 != e1 || r2 != e2 {
  		t.Errorf("expected '%s', '%s' but got '%s', '%s'", e1, e2, r1, r2)
  	}
  }
  ```

  进行测试后的结果如下：

  <img src="https://img-blog.csdnimg.cn/20201026211448266.png#pic_center" alt="在这里插入图片描述" style="zoom: 33%;" />

- `func Test_GlobalFlags(t *testing.T)` 用来测试全局的命令参数

  ```go
  // 测试全局的命令参数
  func Test_GlobalFlags(t *testing.T) {
  	test1.AddCommand(test2)
  	test2.getGflags().StringP("test", "t", "", "")
  	args := []string{"-thupf"}
  	test2.ParseFlags(args)
  	r, _ := test1.getGflags().GetString("test")
  
  	expected := "hupf"
  	if r != expected {
  		t.Errorf("expected '%s', but got '%s'", expected, r)
  	}
  }
  ```

  进行测试后的结果如下：

  <img src="https://img-blog.csdnimg.cn/20201026211602300.png#pic_center" alt="在这里插入图片描述" style="zoom: 33%;" />

- `func Test_LocalFlags(t *testing.T) ` 测试局部的命令参数

  ```go
  // 测试局部的命令参数
  func Test_LocalFlags(t *testing.T) {
  	test1.AddCommand(test2)
  	test2.LocalFlags().StringP("test", "t", "", "")
  	args := []string{"-thupf"}
  	test2.ParseFlags(args)
  	r1, _ := test1.getGflags().GetString("test")
  	e1 := ""
  	r2, _ := test2.LocalFlags().GetString("test")
  	e2 := "hupf"
  
  	if r1 != e1 || r2 != e2 {
  		t.Errorf("expected '%s', '%s' but got '%s', '%s'", e1, e2, r1, r2)
  	}
  }
  ```

  进行测试后的结果如下：

  <img src="https://img-blog.csdnimg.cn/20201026211738911.png#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

- `func Test_Flags(t *testing.T)` 测试所有的命令参数

  ```go
  // 测试所有的命令参数
  func Test_Flags(t *testing.T) {
  	test1.AddCommand(test2)
  	test2.LocalFlags().StringP("local", "l", "", "")
  	test1.getGflags().StringP("global", "g", "", "")
  	args := []string{"-ltestl", "-gtestg"}
  	test2.ParseFlags(args)
  	test1.ParseFlags(args)
  	r1, _ := test2.Flags().GetString("local")
  	e1 := "testl"
  	r2, _ := test2.Flags().GetString("global")
  	e2 := "testg"
  	if r1 != e1 || r2 != e2 {
  		t.Errorf("expected '%s', '%s' but got '%s', '%s'", e1, e2, r1, r2)
  	}
  }
  ```

  进行测试后的结果如下：

  <img src="https://img-blog.csdnimg.cn/2020102621184340.png#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

- `func Test_CommandPath(t *testing.T)` 测试命令路径的正确性

  ```go
  // 测试命令路径的正确性
  func Test_CommandPath(t *testing.T) {
  	test1.AddCommand(test2)
  	r := test2.CommandPath()
  
  	expected := "test1 test2"
  	if r != expected {
  		t.Errorf("expected '%s', but got '%s'", expected, r)
  	}
  }
  ```

  进行测试后的结果如下：

  <img src="https://img-blog.csdnimg.cn/20201026211944284.png#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

## 功能测试

### 基准测试

基准测试 `bench_test` 用来测试函数运行的时间是如何，该函数的实现过程就是首先定义一个命令，然后执行此命令，进行测试：

```go
package mycobra

import (
	"testing"
)

func BenchmarkCommand_Execute(b *testing.B) {
	var test = &Command{
		Use:     "test",
		Short:   "test",
		Long:    "test test",
		Example: "test",
	}
	for i := 0; i < b.N; i++ {
		test.Execute()
	}
}
```

测试的结果如下所示：

<img src="https://img-blog.csdnimg.cn/20201026212416796.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

### 获取使用案例

我设计了一个 `pinfo` 命令，来实现显示自己的个人信息(personal information)，获取该使用案例的命令如下：

`go get github.com/hupf3/pinfo`

或者在相应的 src 目录结构下是用以下命令：

`git clone https://github.com/hupf3/pinfo.git`

`go install`

执行完毕上面的命令后，在 `GOPATH` 路径中的 bin 文件夹下面会多出一个 pinfo 的可执行文件，即为成功

### 使用案例测试

获取完上面的包后，打开命令行，输入 pinfo 测试是否成功安装 pinfo 命令

<img src="https://img-blog.csdnimg.cn/20201026213110118.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

如果出现如上图所示的 pinfo 命令的说明即可证明成功安装 pinfo 命令

然后按照说明可以得知，该命令是获取个人信息的命令，该命令下有三个字命令，分别为：age 获取年龄， id 获取学号， name获取名字

为了实现带参数的命令，我在 name 子命令中定义了三个参数，-f 获取姓，-g 获取名，-a 获取全称

我还实现了 help 相当于是全局的参数，当命令行输入command + help 时会显示该 command 的使用用法，示例如下：

<img src="https://img-blog.csdnimg.cn/20201026213549530.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom: 33%;" />

<img src="https://img-blog.csdnimg.cn/20201026213645207.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

测试完 help 功能后可以进行测试 age 命令是否实现：

<img src="https://img-blog.csdnimg.cn/20201026213807386.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

测试完 age 功能后可以进行测试 id 命令是否实现：

<img src="https://img-blog.csdnimg.cn/20201026213844794.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

测试完 id 功能后可以进行测试 name 各个命令参数是否实现：

<img src="https://img-blog.csdnimg.cn/2020102621390017.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

至此完成了全部的功能测试

## API文档

生成网页版的 API 文档，输入以下的命令：

`godoc -http=:8080`

然后在浏览器中打开 http://127.0.0.1:8080 ，即可访问网页版的 go doc：

<img src="https://img-blog.csdnimg.cn/20201026214726292.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

然后网页搜索 `mycobra` 即可找到我实现的程序包：

<img src="https://img-blog.csdnimg.cn/20201026214827733.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

点开即可查看我实现的 `mycobra` 包中的函数和具体说明以及引用该包的方法：

<img src="https://img-blog.csdnimg.cn/2020102621492556.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

在目录结构下执行以下命令，即可生成线下的 html 文件

`go doc`

`godoc -url="pkg/github.com/hupf3/mycobra" > API.html`

我将该文档也保存在了 github 仓库中方便检查

## 总结

通过本次实验提升了自己阅读程序包的能力，并且能够根据已有的程序包提取有用的信息，改善程序包，变成自己的程序包。在本次实验也注意到了一个易错点，就是双大括号引用下的函数，函数名的首字母一定要大写

