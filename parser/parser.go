package parser

import (
  "monkey/ast"
  "monkey/lexer"
  "monkey/token"
)

type Parser struct {
  lexer *lexer.Lexer
  currentToken token.Token
  peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
  p := &Parser{lexer: lexer}

  // Read two tokens, so currentToken and peekToken are both set
  p.nextToken()
  p.nextToken()

  return p
}

func (parser *Parser) nextToken() {
  parser.currentToken = parser.peekToken
  parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
  program := &ast.Program{}
  program.Statements = []ast.Statement{}

  for !parser.currentTokenIs(token.EOF) {
    statement := parser.parseStatement()
    if statement != nil {
      program.Statements = append(program.Statements, statement)
    }
    parser.nextToken()
  }

  return program
}

func (parser *Parser) parseStatement() ast.Statement {
  switch parser.currentToken.Type {
  case token.LET:
    return parser.parseLetStatement()
  default:
    return nil
  }
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
  statement := &ast.LetStatement{Token: parser.currentToken}

  if !parser.expectPeek(token.IDENT) {
    return nil
  }

  statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

  if !parser.expectPeek(token.ASSIGN) {
    return nil
  }

  // We're skipping the epxrssions until we encounter a semicolon
  for !parser.currentTokenIs(token.SEMICOLON) {
    parser.nextToken()
  }

  return statement
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
  return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
  return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
  if parser.peekTokenIs(tokenType) {
    parser.nextToken()
    return true
  }

  return false
}
