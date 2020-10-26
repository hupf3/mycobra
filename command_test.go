package mycobra

import (
	"testing"
)

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

// 测试命令带参数
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

// 测试局部的命令参数
func Test_LocalFlags(t *testing.T) {
	test1.AddCommand(test2)
	test2.LocalFlags().StringP("test", "t", "", "")
	args := []string{"-thupf"}
	test2.ParseFlags(args)
	r1, _ := test1.getGflags().GetString("test")
	e1 := "hupf"
	r2, _ := test2.LocalFlags().GetString("test")
	e2 := ""

	if r1 != e1 || r2 != e2 {
		t.Errorf("expected '%s', '%s' but got '%s', '%s'", e1, e2, r1, r2)
	}
}

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

// 测试命令路径的正确性
func Test_CommandPath(t *testing.T) {
	test1.AddCommand(test2)
	r := test2.CommandPath()

	expected := "test1 test2"
	if r != expected {
		t.Errorf("expected '%s', but got '%s'", expected, r)
	}
}
