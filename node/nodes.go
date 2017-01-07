package node

import (
	"github.com/SteveZhangBit/leiogo-css/parser"
	"golang.org/x/net/html"
	"strings"
)

type Node html.Node

func Find(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.ElementNode {
			if c.IsMatch(query) {
				nodes = append(nodes, c)
			}
			nodes = append(nodes, Find(c, query)...)
		}
	}
	return nodes
}

func Child(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.ElementNode && (len(query.Seq) == 0 || c.IsMatch(query)) {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

func Not(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.ElementNode && !c.IsMatch(query) {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

func Next(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for s := (*Node)(n.NextSibling); s != nil; s = (*Node)(s.NextSibling) {
		if s.Type == html.ElementNode {
			if len(query.Seq) == 0 || s.IsMatch(query) {
				return append(nodes, s)
			} else {
				return nodes
			}
		}
	}
	return nodes
}

func NextAll(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for s := (*Node)(n.NextSibling); s != nil; s = (*Node)(s.NextSibling) {
		if s.Type == html.ElementNode && (len(query.Seq) == 0 || s.IsMatch(query)) {
			nodes = append(nodes, s)
		}
	}
	return nodes
}

func Prev(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for s := (*Node)(n.PrevSibling); s != nil; s = (*Node)(s.PrevSibling) {
		if s.Type == html.ElementNode {
			if len(query.Seq) == 0 || s.IsMatch(query) {
				return append(nodes, s)
			} else {
				return nodes
			}
		}
	}
	return nodes
}

func PrevAll(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for s := (*Node)(n.PrevSibling); s != nil; s = (*Node)(s.PrevSibling) {
		if s.Type == html.ElementNode && (len(query.Seq) == 0 || s.IsMatch(query)) {
			nodes = append(nodes, s)
		}
	}
	return nodes
}

func Parent(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	if s := (*Node)(n.Parent); s != nil {
		if s.Type == html.ElementNode && (len(query.Seq) == 0 || s.IsMatch(query)) {
			return append(nodes, s)
		}
	}
	return nodes
}

func Parents(n *Node, query parser.Element) []*Node {
	nodes := []*Node{}
	for s := (*Node)(n.Parent); s != nil; s = (*Node)(s.Parent) {
		if s.Type == html.ElementNode && (len(query.Seq) == 0 || s.IsMatch(query)) {
			nodes = append(nodes, s)
		}
	}
	return nodes
}

func (n *Node) IsMatch(query parser.Element) (isMatch bool) {
	for _, ast := range query.Seq {
		switch x := ast.(type) {
		case parser.Tag:
			if isMatch = (x.Name == "*" || n.Data == x.Name); !isMatch {
				break
			}
		case parser.Id:
			if isMatch = n.GetId() == x.Name; !isMatch {
				break
			}
		case parser.Class:
			isMatch = false
			for _, c := range n.GetCls() {
				if c == x.Name {
					isMatch = true
					break
				}
			}
		case parser.Attr:
			switch x.Type {
			case "":
				if isMatch = n.HasAttr(x.Name); !isMatch {
					break
				}
			case "=":
				if isMatch = x.Value == n.GetAttr(x.Name); !isMatch {
					break
				}
			case "^=":
				if isMatch = strings.HasPrefix(n.GetAttr(x.Name), x.Value); !isMatch {
					break
				}
			case "$=":
				if isMatch = strings.HasSuffix(n.GetAttr(x.Name), x.Value); !isMatch {
					break
				}
			case "*=":
				if isMatch = strings.Contains(n.GetAttr(x.Name), x.Value); !isMatch {
					break
				}
			}
		}
	}
	return
}

func (n *Node) GetId() string {
	for _, attr := range n.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}
	return ""
}

func (n *Node) GetCls() []string {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return strings.Split(attr.Val, " ")
		}
	}
	return nil
}

func (n *Node) HasAttr(name string) bool {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return true
		}
	}
	return false
}

func (n *Node) GetAttr(name string) string {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func (n *Node) Text() string {
	s := ""
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.TextNode {
			s += c.Data
		}
	}
	return s
}
