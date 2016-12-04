package parser

import (
	"fmt"
	"testing"
)

func test(t *testing.T, str string) {
	p := NewParser(str)
	if ast, err := p.Parse(); err != nil {
		t.Error(err.Error())
	} else {
		fmt.Println(PrintVisitor(ast))
	}
}

func Test1(t *testing.T) {
	test(t, "a")
}

func Test2(t *testing.T) {
	test(t, "#id")
}

func Test3(t *testing.T) {
	test(t, ".cls")
}

func Test4(t *testing.T) {
	test(t, "[attr]")
}

func Test5(t *testing.T) {
	test(t, "a#id.cls[attr]")
}

func Test6(t *testing.T) {
	test(t, "div a[attr=123]")
}

func Test7(t *testing.T) {
	test(t, "div > img")
}

func Test8(t *testing.T) {
	test(t, "div a.cls[attr^=123.45] > img")
}

func Test9(t *testing.T) {
	test(t, "div + img[src*=htt://www.baidu.com] ~ [attr$=.png]")
}

func Test10(t *testing.T) {
	test(t, "a, img[src=\"abc\"], div h3.cls[attr='abc.123']")
}
