package lexer

import (
	"fmt"
	"testing"
)

func compare(a []Token, b []Token) (bool, string) {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return false, fmt.Sprintf("Get %v, need %v", a[i], b[i])
		}
	}
	return true, ""
}

func ParseTest(input string, output ...Token) (bool, string) {
	lexer := NewLexer(input)
	tokens := make([]Token, 0)
	for {
		if t, err := lexer.NextToken(); err != nil {
			return false, err.Error()
		} else {
			tokens = append(tokens, t)
			if t.Type == EOF {
				return compare(tokens, output)
			}
		}
	}
}

func Test1(t *testing.T) {
	if ok, err :=
		ParseTest(
			"a",
			Token{Type: Identifier, Value: "a"},
			Token{Type: EOF}); !ok {
		t.Error(err)
	}
}

func Test2(t *testing.T) {
	if ok, err :=
		ParseTest(
			"a.class img",
			Token{Type: Identifier, Value: "a"},
			Token{Type: Dot},
			Token{Type: Identifier, Value: "class"},
			Token{Type: Blank},
			Token{Type: Identifier, Value: "img"},
			Token{Type: EOF}); !ok {
		t.Error(err)
	}
}

func Test3(t *testing.T) {
	if ok, err :=
		ParseTest(
			"a#id  >  img.cls1.cls2",
			Token{Type: Identifier, Value: "a"},
			Token{Type: Sharp},
			Token{Type: Identifier, Value: "id"},
			Token{Type: Blank},
			Token{Type: Greater},
			Token{Type: Blank},
			Token{Type: Identifier, Value: "img"},
			Token{Type: Dot},
			Token{Type: Identifier, Value: "cls1"},
			Token{Type: Dot},
			Token{Type: Identifier, Value: "cls2"},
			Token{Type: EOF}); !ok {
		t.Error(err)
	}
}

func Test4(t *testing.T) {
	if ok, err :=
		ParseTest(
			"[attr] ~  .cls[attr=123.45]",
			Token{Type: LeftBracket},
			Token{Type: Identifier, Value: "attr"},
			Token{Type: RightBracket},
			Token{Type: Blank},
			Token{Type: Wave},
			Token{Type: Blank},
			Token{Type: Dot},
			Token{Type: Identifier, Value: "cls"},
			Token{Type: LeftBracket},
			Token{Type: Identifier, Value: "attr"},
			Token{Type: Assign},
			Token{Type: Literal, Value: "123.45"},
			Token{Type: RightBracket},
			Token{Type: EOF}); !ok {
		t.Error(err)
	}
}

func Test5(t *testing.T) {
	if ok, err :=
		ParseTest(
			"a[attr='123'] + img[src=http://www.baidu.com], img[href=\"abc\"]",
			Token{Type: Identifier, Value: "a"},
			Token{Type: LeftBracket},
			Token{Type: Identifier, Value: "attr"},
			Token{Type: Assign},
			Token{Type: String, Value: "'123'"},
			Token{Type: RightBracket},
			Token{Type: Blank},
			Token{Type: Plus},
			Token{Type: Blank},
			Token{Type: Identifier, Value: "img"},
			Token{Type: LeftBracket},
			Token{Type: Identifier, Value: "src"},
			Token{Type: Assign},
			Token{Type: Literal, Value: "http://www.baidu.com"},
			Token{Type: RightBracket},
			Token{Type: Comma},
			Token{Type: Blank},
			Token{Type: Identifier, Value: "img"},
			Token{Type: LeftBracket},
			Token{Type: Identifier, Value: "href"},
			Token{Type: Assign},
			Token{Type: String, Value: "\"abc\""},
			Token{Type: RightBracket},
			Token{Type: EOF}); !ok {
		t.Error(err)
	}
}
