package parser

import (
  "testing"
  "monkey/ast"
  "monkey/lexer"
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
    t.Fatalf("ParseProgram() returned nil")
  }

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements count does not equal 3, received %d statements", len(program.Statements))
  }

  tests := []struct {
    expectedIdentifier string
  }{
    {"x"},
    {"y"},
    {"foobar"},
  }

  for i, tt := range tests {
    statement := program.Statements[i]
    if !testLetStatement(t, statement, tt.expectedIdentifier) {
      return
    }
  }
}

func TestReturnStatements(t *testing.T) {
  input := `return 5; return 10; return 993322;`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contains 3 statements. Received %q",
      len(program.Statements))
  }

  for _, statement := range program.Statements {
    returnStatement, ok := statement.(*ast.ReturnStatement)
    
    if !ok {
      t.Errorf("statement not of type ReturnStatement, got=%T", statement)
      continue
    }
    
    if returnStatement.TokenLiteral() != "return" {
      t.Errorf("returnStatement.TokenLiteral is not 'return', got %q", returnStatement.TokenLiteral())
    }
  }
}  

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
  if s.TokenLiteral() != "let" {
    t.Errorf("s.TokenLiteral() not a 'let'. received %q", s.TokenLiteral())
  }

  letStatement, ok := s.(*ast.LetStatement)

  if !ok {
    t.Errorf("s not *ast.LetStatement. got=%T", s)
    return false
  }

  if letStatement.Name.Value != name {
    t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStatement.Name.Value) 
    return false
  }

  if letStatement.Name.TokenLiteral() != name { 
    t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
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
    t.Errorf("Parser error: %q", msg)
  }

  t.FailNow()
}
