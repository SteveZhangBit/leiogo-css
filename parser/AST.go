package parser

type AST interface{}

type Selector struct {
	Seq []AST
}

type Exp struct {
	E  AST
	F  AST
	Op string
}

type Element struct {
	Seq []AST
}

type Tag struct {
	Name string
}

type Id struct {
	Name string
}

type Class struct {
	Name string
}

type Attr struct {
	Name  string
	Value string
	Type  string
}

type ASTBuilder struct {
	stack []AST
	count int
	err   error
}

func (b *ASTBuilder) push(val AST) {
	if b.err != nil {
		return
	}
	b.stack = append(b.stack, val)
}

func (b *ASTBuilder) pop() AST {
	if len(b.stack) == 0 {
		return nil
	} else {
		length := len(b.stack)
		val := b.stack[length-1]
		b.stack = b.stack[:length-1]
		return val
	}
}

func (b *ASTBuilder) tag(name string) {
	if b.err != nil {
		return
	}
	b.push(Tag{Name: name})
	b.count++
}

func (b *ASTBuilder) id(name string) {
	if b.err != nil {
		return
	}
	b.push(Id{Name: name})
	b.count++
}

func (b *ASTBuilder) class(name string) {
	if b.err != nil {
		return
	}
	b.push(Class{Name: name})
	b.count++
}

func (b *ASTBuilder) attr(t string) {
	if b.err != nil {
		return
	}
	if t == "" {
		b.push(Attr{Name: b.pop().(string)})
	} else {
		val := b.pop().(string)
		name := b.pop().(string)
		b.push(Attr{Name: name, Value: val, Type: t})
	}
	b.count++
}

func (b *ASTBuilder) element() {
	if b.err != nil {
		return
	}
	el := Element{Seq: make([]AST, b.count)}
	for b.count--; b.count >= 0; b.count-- {
		el.Seq[b.count] = b.pop()
	}
	b.count = 0
	b.push(el)
}

func (b *ASTBuilder) exp(op string) {
	if b.err != nil {
		return
	}
	F := b.pop()
	E := b.pop()
	b.push(Exp{E: E, F: F, Op: op})
}

func (b *ASTBuilder) selector() {
	if b.err != nil {
		return
	}
	sel := Selector{Seq: make([]AST, len(b.stack))}
	for i := len(b.stack) - 1; i >= 0; i-- {
		sel.Seq[i] = b.pop()
	}
	b.push(sel)
}

func (b *ASTBuilder) build() AST {
	if b.err != nil {
		return nil
	}
	return b.pop()
}
