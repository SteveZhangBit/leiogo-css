package lexer

import (
	"errors"
	"fmt"
	"unicode"
)

type status int

const (
	status0 status = iota
	status1
)

type Lexer struct {
	str    []rune
	i      int
	length int
	status status
}

var tokenNames = []string{
	"",             // error token
	"Comma",        // ,
	"LeftBracket",  // [
	"RightBracket", // ]
	"Greater",      // >
	"Plus",         // +
	"Wave",         // ~
	"Sharp",        // #
	"Dot",          // .
	"Up",           // ^
	"Dollar",       // $
	"Star",         // *
	"Assign",       // =
	"Identifier",   // [a-zA-Z_]\w*
	"Literal",      // 123, 1.0, /path/
	"String",       // '...', "..."
	"Blank",
	"EOF",
}

type TokenType int

func (t TokenType) String() string {
	return tokenNames[t]
}

const (
	Comma        TokenType = 1 + iota // ,
	LeftBracket                       // [
	RightBracket                      // ]
	Greater                           // >
	Plus                              // +
	Wave                              // ~
	Sharp                             // #
	Dot                               // .
	Up                                // ^
	Dollar                            // $
	Star                              // *
	Assign                            // =
	Identifier                        // [a-zA-Z_]\w*
	Literal                           // 123, 1.0, /path/
	String                            // '...', "..."
	Blank
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("<%s, %s>", t.Type, t.Value)
}

func NewLexer(str string) *Lexer {
	return &Lexer{
		str:    []rune(str),
		length: len(str),
		i:      -1,
		status: status0,
	}
}

func (l *Lexer) getRune() (rune, error) {
	l.i++
	if l.i < l.length {
		return l.str[l.i], nil
	}
	return 0, errors.New("EOF")
}

func (l *Lexer) NextToken() (token Token, err error) {
	var c rune

	if c, err = l.getRune(); err != nil {
		return Token{Type: EOF}, nil
	}

	if unicode.IsSpace(c) {
		for {
			if c, err = l.getRune(); err != nil || !unicode.IsSpace(c) {
				l.i--
				return Token{Type: Blank}, nil
			}
		}
	}

	switch l.status {
	case status0:
		switch c {
		case ',':
			token = Token{Type: Comma}
		case '[':
			token = Token{Type: LeftBracket}
		case ']':
			token = Token{Type: RightBracket}
		case '>':
			token = Token{Type: Greater}
		case '+':
			token = Token{Type: Plus}
		case '~':
			token = Token{Type: Wave}
		case '#':
			token = Token{Type: Sharp}
		case '.':
			token = Token{Type: Dot}
		case '^':
			token = Token{Type: Up}
		case '$':
			token = Token{Type: Dollar}
		case '*':
			token = Token{Type: Star}
		case '=':
			l.status = status1
			token = Token{Type: Assign}
		default:
			// Identifier
			if unicode.IsLetter(c) || c == '_' {
				id := []rune{c}
				for {
					if c, err = l.getRune(); err == nil &&
						(unicode.IsLetter(c) || c == '_' || c == '-' || unicode.IsDigit(c)) {
						id = append(id, c)
					} else {
						l.i-- // push back a rune
						return Token{Type: Identifier, Value: string(id)}, nil
					}
				}
			} else {
				return Token{}, errors.New("Unexpected rune: " + string(c))
			}
		}
	case status1:
		switch c {
		case ']':
			l.status = status0
			token = Token{Type: RightBracket}
		// String
		case '\'':
			s := []rune{c}
			for {
				if c, err = l.getRune(); err != nil {
					return Token{}, errors.New("Unclosed string: " + string(s))
				} else if c == '\'' {
					s = append(s, c)
					return Token{Type: String, Value: string(s)}, nil
				} else {
					s = append(s, c)
				}
			}
		case '"':
			s := []rune{c}
			for {
				if c, err = l.getRune(); err != nil {
					return Token{}, errors.New("Unclosed string: " + string(s))
				} else if c == '"' {
					s = append(s, c)
					return Token{Type: String, Value: string(s)}, nil
				} else {
					s = append(s, c)
				}
			}
		default:
			s := []rune{c}
			for {
				if c, err = l.getRune(); err != nil || c == '"' || c == '\'' || c == ']' {
					l.i--
					return Token{Type: Literal, Value: string(s)}, nil
				} else {
					s = append(s, c)
				}
			}
		}
	default:
		token, err = Token{}, errors.New("Unknown status")
	}
	return
}
