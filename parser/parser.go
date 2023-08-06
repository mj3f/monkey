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
  return nil
}
