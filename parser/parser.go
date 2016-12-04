package parser

import (
	"errors"
	"fmt"

	"github.com/SteveZhangBit/css/lexer"
)

type Parser struct {
	lexer     *lexer.Lexer
	lookahead lexer.Token
	Builder   ASTBuilder
}

func NewParser(str string) *Parser {
	p := Parser{lexer: lexer.NewLexer(str)}
	p.match(lexer.Token{}.Type)
	return &p
}

func (p *Parser) Parse() (AST, error) {
	p.entry()
	return p.Builder.build(), p.Builder.err
}

func (p *Parser) match(t lexer.TokenType) {
	if p.Builder.err != nil {
		return
	}

	if p.lookahead.Type == t {
		p.lookahead, p.Builder.err = p.lexer.NextToken()
	} else {
		p.Builder.err = errors.New(fmt.Sprintf("Need %s, Get %s", t, p.lookahead))
	}
}

func (p *Parser) entry() {
	p.selector()
	p.Builder.selector()
}

func (p *Parser) selector() {
	p.exp()
	p.selector_()
}

func (p *Parser) selector_() {
	switch p.lookahead.Type {
	case lexer.Comma:
		p.match(lexer.Comma)
		p.space()
		p.selector()
	case lexer.EOF:
		return
	default:
		p.err(lexer.Comma, lexer.EOF)
	}
}

func (p *Parser) space() {
	if p.lookahead.Type == lexer.Blank {
		p.match(lexer.Blank)
	}
}

func (p *Parser) exp() {
	p.element()
	p.exp_()
}

func (p *Parser) exp_() {
	switch p.lookahead.Type {
	case lexer.Blank:
		p.match(lexer.Blank)
		p.expCombine()
	case lexer.Greater:
		p.child()
	case lexer.Plus:
		p.imPrecedent()
	case lexer.Wave:
		p.precedent()
	case lexer.Comma, lexer.EOF:
		return
	default:
		p.err(lexer.Blank, lexer.Greater, lexer.Plus, lexer.Wave)
	}
}

func (p *Parser) expCombine() {
	switch p.lookahead.Type {
	case lexer.Greater:
		p.child()
	case lexer.Plus:
		p.imPrecedent()
	case lexer.Wave:
		p.precedent()
	default:
		p.descendant()
	}
}

func (p *Parser) descendant() {
	p.element()
	p.Builder.exp(" ")
	p.exp_()
}

func (p *Parser) child() {
	p.match(lexer.Greater)
	p.space()
	p.element()
	p.Builder.exp(">")
	p.exp_()
}

func (p *Parser) imPrecedent() {
	p.match(lexer.Plus)
	p.space()
	p.element()
	p.Builder.exp("+")
	p.exp_()
}

func (p *Parser) precedent() {
	p.match(lexer.Wave)
	p.space()
	p.element()
	p.Builder.exp("~")
	p.exp_()
}

func (p *Parser) element() {
	if p.lookahead.Type == lexer.Identifier || p.lookahead.Type == lexer.Star {
		p.tag()
	}
	p.adjunct()
	p.Builder.element()
}

func (p *Parser) adjunct() {
	switch p.lookahead.Type {
	case lexer.Sharp:
		p.id()
		p.adjunct()
	case lexer.Dot:
		p.class()
		p.adjunct()
	case lexer.LeftBracket:
		p.attr()
		p.adjunct()
	case lexer.Blank, lexer.Greater, lexer.Plus, lexer.Wave, lexer.Comma, lexer.EOF:
		return
	default:
		p.err(lexer.Sharp, lexer.Dot, lexer.LeftBracket)
	}
}

func (p *Parser) id() {
	p.match(lexer.Sharp)
	t := p.lookahead
	p.match(lexer.Identifier)
	p.Builder.id(t.Value)
}

func (p *Parser) class() {
	p.match(lexer.Dot)
	t := p.lookahead
	p.match(lexer.Identifier)
	p.Builder.class(t.Value)
}

func (p *Parser) attr() {
	p.match(lexer.LeftBracket)
	t := p.lookahead
	p.match(lexer.Identifier)
	p.Builder.push(t.Value)

	switch p.lookahead.Type {
	case lexer.RightBracket:
		p.match(lexer.RightBracket)
		p.Builder.attr("")
	case lexer.Assign:
		p.match(lexer.Assign)
		p.literal()
		p.match(lexer.RightBracket)
		p.Builder.attr("=")
	case lexer.Up:
		p.match(lexer.Up)
		p.match(lexer.Assign)
		p.literal()
		p.match(lexer.RightBracket)
		p.Builder.attr("^=")
	case lexer.Dollar:
		p.match(lexer.Dollar)
		p.match(lexer.Assign)
		p.literal()
		p.match(lexer.RightBracket)
		p.Builder.attr("$=")
	case lexer.Star:
		p.match(lexer.Star)
		p.match(lexer.Assign)
		p.literal()
		p.match(lexer.RightBracket)
		p.Builder.attr("*=")
	default:
		p.err(lexer.RightBracket, lexer.Assign, lexer.Up, lexer.Dollar, lexer.Star)
	}
}

func (p *Parser) tag() {
	switch t := p.lookahead; t.Type {
	case lexer.Identifier:
		p.match(lexer.Identifier)
		p.Builder.tag(t.Value)
	case lexer.Star:
		p.match(lexer.Star)
		p.Builder.tag("*")
	default:
		p.err(lexer.Identifier, lexer.Star)
	}
}

func (p *Parser) literal() {
	switch t := p.lookahead; t.Type {
	case lexer.String:
		p.match(lexer.String)
		p.Builder.push(t.Value[1 : len(t.Value)-1])
	case lexer.Literal:
		p.match(lexer.Literal)
		p.Builder.push(t.Value)
	default:
		p.err(lexer.String, lexer.Literal)
	}
}

func (p *Parser) err(need ...lexer.TokenType) {
	p.Builder.err = errors.New(fmt.Sprintf("Need %s, get %s", need, p.lookahead))
}
