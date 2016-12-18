package selector

import (
	"strings"

	"github.com/SteveZhangBit/css/node"
	"github.com/SteveZhangBit/css/parser"
	"golang.org/x/net/html"
)

type Elements struct {
	Nodes []*node.Node
	Err   error
}

func Parse(body string) *Elements {
	if doc, err := html.Parse(strings.NewReader(body)); err != nil {
		return &Elements{Err: err}
	} else {
		return &Elements{Nodes: []*node.Node{(*node.Node)(doc)}}
	}
}

func (e *Elements) selectorHelper(str string, f func(ast parser.AST) *Elements) *Elements {
	if e.Err != nil {
		return e
	}
	if ast, err := parser.NewParser(str).Parse(); err != nil {
		e.Err = err
		return e
	} else {
		return f(ast)
	}
}

func (e *Elements) Find(str string) *Elements {
	return e.selectorHelper(str, e.find)
}

func (e *Elements) Child(str string) *Elements {
	return e.selectorHelper(str, e.child)
}

func (e *Elements) Not(str string) *Elements {
	return e.selectorHelper(str, e.not)
}

func (e *Elements) Next(str string) *Elements {
	return e.selectorHelper(str, e.next)
}

func (e *Elements) NextAll(str string) *Elements {
	return e.selectorHelper(str, e.nextAll)
}

func (e *Elements) Prev(str string) *Elements {
	return e.selectorHelper(str, e.prev)
}

func (e *Elements) PrevAll(str string) *Elements {
	return e.selectorHelper(str, e.prevAll)
}

func (e *Elements) Parent(str string) *Elements {
	return e.selectorHelper(str, e.parent)
}

func (e *Elements) Parents(str string) *Elements {
	return e.selectorHelper(str, e.parents)
}

func (e *Elements) Text() string {
	if len(e.Nodes) == 0 {
		return ""
	} else {
		return e.Nodes[0].Text()
	}
}

func (e *Elements) Texts() []string {
	text := []string{}
	for _, n := range e.Nodes {
		text = append(text, n.Text())
	}
	return text
}

func (e *Elements) Attr(name string) string {
	if len(e.Nodes) == 0 {
		return ""
	} else {
		return e.Nodes[0].GetAttr(name)
	}
}

func (e *Elements) Attrs(name string) []string {
	attr := []string{}
	for _, n := range e.Nodes {
		attr = append(attr, n.GetAttr(name))
	}
	return attr
}

func (e *Elements) find(ast parser.AST) *Elements {
	nodes := []*node.Node{}
	switch x := ast.(type) {
	case parser.Selector:
		for _, exp := range x.Seq {
			nodes = append(nodes, e.find(exp).Nodes...)
		}
	case parser.Exp:
		E := e.find(x.E)
		switch x.Op {
		case " ":
			return E.find(x.F)
		case ">":
			return E.child(x.F)
		case "+":
			return E.next(x.F)
		case "~":
			return E.nextAll(x.F)
		}
	case parser.Element:
		for _, n := range e.Nodes {
			nodes = append(nodes, node.Find(n, x)...)
		}
	}
	return &Elements{Nodes: nodes}
}

func (e *Elements) selectorHelper2(ast parser.AST, f func(n *node.Node, query parser.Element) []*node.Node) *Elements {
	nodes := []*node.Node{}
	switch x := ast.(type) {
	case parser.Selector:
		for _, exp := range x.Seq {
			nodes = append(nodes, e.selectorHelper2(exp, f).Nodes...)
		}
	case parser.Element:
		for _, n := range e.Nodes {
			nodes = append(nodes, f(n, x)...)
		}
	}
	return &Elements{Nodes: nodes}
}

func (e *Elements) child(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Child)
}

func (e *Elements) not(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Not)
}

func (e *Elements) next(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Next)
}

func (e *Elements) nextAll(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.NextAll)
}

func (e *Elements) prev(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Prev)
}

func (e *Elements) prevAll(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.PrevAll)
}

func (e *Elements) parent(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Parent)
}

func (e *Elements) parents(ast parser.AST) *Elements {
	return e.selectorHelper2(ast, node.Parents)
}
