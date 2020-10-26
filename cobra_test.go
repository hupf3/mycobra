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
