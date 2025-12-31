package parser

import (
	"fmt"
	"testing"
	"github.com/aldosierra/monkey/ast"
	"github.com/aldosierra/monkey/lexer"
)


func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if statelen := len(program.Statements); statelen != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			statelen)
	}

	tests := []struct {
		expectedIndetifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIndetifier) {
			return
		}
	}
}


func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if statelen := len(program.Statements); statelen != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			statelen)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnType := returnStmt.TokenLiteral(); returnType != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return' got=%q",
				returnType)
		}
	}

	// tests := []struct{
	// 	input         string
	// 	expectedValue any
	// }{
	// 	{ "return 5;",     5 },
	// 	{ "return 10;",    10 },
	// 	{ "return 993322", 993322 },
	// 	{ "return a",      "a" },
	// }
	//
	// for _, tt := range tests {
	// 	l := lexer.New(tt.input)
	// 	p := New(l)
	//
	// 	program := p.ParseProgram()
	// 	checkParserErrors(t, p)
	//
	// 	if statelen := len(program.Statements); statelen != 1 {
	// 		t.Fatalf("program.Statements does not contain 1 statement. got=%d", statelen)
	// 	}
	//
	// 	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	//
	// 	if !ok {
	// 		t.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T",
	// 			program.Statements[0])
	// 	}
	//
	// 	if tk := stmt.TokenLiteral(); tk != "return" {
	// 		t.Fatalf("stmt.TokenLiteral is not 'return'. got=%q", tk)
	// 	}
	//
	// 	if !testLiteralExpression(t, stmt.ReturnValue, tt.expectedValue) {
	// 		return
	// 	}
	// }
}


func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if statelen := len(program.Statements); statelen != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			statelen)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}


func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if statelen := len(program.Statements); statelen != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			statelen)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}


func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct{
		input        string
		operator     string
		integerValue int64
	}{
		{ "!5", "!", 5 },
		{ "-15", "-", 15 },
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if statelen := len(program.Statements); statelen != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, statelen)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt exp is not ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}


func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct{
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{ "5 + 5;", 5, "+", 5 },
		{ "5 - 5;", 5, "-", 5 },
		{ "5 * 5;", 5, "*", 5 },
		{ "5 / 5;", 5, "/", 5 },
		{ "5 > 5;", 5, ">", 5 },
		{ "5 < 5;", 5, "<", 5 },
		{ "5 == 5;", 5, "==", 5 },
		{ "5 != 5;", 5, "!=", 5 },
		{ "a != a;", "a", "!=", "a" },
		{ "a != b;", "a", "!=", "b" },
		{ "a * b;", "a", "*", "b" },
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if statelen := len(program.Statements); statelen != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, statelen)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator,
			tt.rightValue) {
			return
		}
	}
}


func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if actual := program.String(); actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}


func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if literal := s.TokenLiteral(); literal != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", literal)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}


func testInfixExpression(t *testing.T, exp ast.Expression, left any,
	operator string, right any) bool {

	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}


func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(value))
	case int64:
		return testIntegerLiteral(t, exp, value)
	case string:
		return testIdentifier(t, exp, value)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}


func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if tk := integ.TokenLiteral(); tk != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, tk)
		return false
	}

	return true
}


func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if tk := ident.TokenLiteral(); tk != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, tk)
		return false
	}

	return true
}


func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
