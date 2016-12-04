package parser

import (
	"fmt"
)

func PrintVisitor(ast AST) (str string) {
	switch x := ast.(type) {
	case Selector:
		str = ""
		for _, exp := range x.Seq {
			fmt.Sprintf("%T\n", exp)
			str += PrintVisitor(exp) + ", "
		}
		str = str[:len(str)-2]
		return
	case Exp:
		str = fmt.Sprintf("[%s%s%s]", PrintVisitor(x.E), x.Op, PrintVisitor(x.F))
		return
	case Element:
		str = "["
		for _, el := range x.Seq {
			str += PrintVisitor(el) + " "
		}
		str = str[:len(str)-1] + "]"
		return
	case Tag:
		return x.Name
	case Id:
		return "#" + x.Name
	case Class:
		return "." + x.Name
	case Attr:
		return fmt.Sprintf("[%s%s%s]", x.Name, x.Type, x.Value)
	default:
		panic(fmt.Sprintf("Error type: %T", x))
	}
}
