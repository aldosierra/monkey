package ast


import (
	"testing"
	"github.com/aldosierra/monkey/token"
)


func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if str := program.String(); str != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", str)
	}
}
