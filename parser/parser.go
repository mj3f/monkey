package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
  prefixParseFn func() ast.Expression
  infixParseFn func(ast.Expression) ast.Expression // param is left side of the infix operator
)

const ( // ORDER OF PRECENDENCE FOR OPERANDS
  _ int = iota // gives the values below incrementing values, 0, 1, 2, etc.
  LOWEST // (1)
  EQUALS // == (2)
  LESSGREATER // > or < (3)
  SUM // + (4)
  PRODUCT // * (5)
  PREFIX // -X OR !X (6)
  CALL // myFn(X) (7)
  INDEX // array[index] (8)
) // e.g. PDOCUT (*) has higher order precedence than EQUALS (==)

var precedences = map[token.TokenType]int {
  token.EQ:           EQUALS,
  token.NOT_EQ:       EQUALS,
  token.LT:           LESSGREATER,
  token.GT:           LESSGREATER,
  token.PLUS:         SUM,
  token.MINUS:        SUM,
  token.SLASH:        PRODUCT,
  token.ASTERISK:     PRODUCT,
  token.LPAREN:       CALL,
  token.LBRACKET:     INDEX,
}

type Parser struct {
  lexer *lexer.Lexer
  currentToken token.Token
  peekToken token.Token
  errors []string
  prefixParseFns map[token.TokenType]prefixParseFn
  infixParseFns map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
  p := &Parser{
    lexer: lexer,
    errors: []string{},
  }

  // Read two tokens, so currentToken and peekToken are both set
  p.nextToken()
  p.nextToken()

  p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
  p.registerPrefix(token.IDENT, p.parseIdentifier)
  p.registerPrefix(token.INT, p.parseIntegerLiteral)
  p.registerPrefix(token.BANG, p.parsePrefixExpression)
  p.registerPrefix(token.MINUS, p.parsePrefixExpression)
  p.registerPrefix(token.TRUE, p.parseBoolean)
  p.registerPrefix(token.FALSE, p.parseBoolean)
  p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
  p.registerPrefix(token.IF, p.parseIfExpression)
  p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
  p.registerPrefix(token.STRING, p.parseStringLiteral)
  p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
  p.registerPrefix(token.LBRACE, p.parseHashLiteral)

  p.infixParseFns = make(map[token.TokenType]infixParseFn)
  p.registerInfix(token.PLUS, p.parseInfixExpression)
  p.registerInfix(token.MINUS, p.parseInfixExpression) 
  p.registerInfix(token.SLASH, p.parseInfixExpression)
  p.registerInfix(token.ASTERISK, p.parseInfixExpression)
  p.registerInfix(token.EQ, p.parseInfixExpression)
  p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
  p.registerInfix(token.LT, p.parseInfixExpression)
  p.registerInfix(token.GT, p.parseInfixExpression)
  p.registerInfix(token.LPAREN, p.parseCallExpression)
  p.registerInfix(token.LBRACKET, p.parseIndexExpression)

  return p
}

func (parser *Parser) Errors() []string {
  return parser.errors
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
  case token.RETURN:
    return parser.parseReturnStatement()
  default:
    return parser.parseExpressionStatement()
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

  parser.nextToken()

  statement.Value = parser.parseExpression(LOWEST)

  for parser.peekTokenIs(token.SEMICOLON) {
    parser.nextToken()
  }

  return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
  statement := &ast.ReturnStatement{Token: parser.currentToken}
  parser.nextToken()

  statement.ReturnValue = parser.parseExpression(LOWEST)

  if parser.peekTokenIs(token.SEMICOLON) {
    parser.nextToken()
  }

  return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
  statement := &ast.ExpressionStatement{Token: parser.currentToken}
  statement.Expression = parser.parseExpression(LOWEST)

  if parser.peekTokenIs(token.SEMICOLON) {
    parser.nextToken()
  }

  return statement

}

func (parser *Parser) parseExpression(precendence int) ast.Expression {
  prefix := parser.prefixParseFns[parser.currentToken.Type]

  if prefix == nil {
    parser.noPrefixParseFnError(parser.currentToken.Type)
    return nil
  }

  leftExpr := prefix()

  for !parser.peekTokenIs(token.SEMICOLON) && precendence < parser.peekPrecedence() {
    infix := parser.infixParseFns[parser.peekToken.Type]
    if infix == nil {
      return leftExpr
    }

    parser.nextToken()

    leftExpr = infix(leftExpr)
  }

  return leftExpr
}

// Prefix parser functions

func (parser *Parser) parseIntegerLiteral() ast.Expression {
  literal := &ast.IntegerLiteral{Token: parser.currentToken}

  value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

  if err != nil {
    msg := fmt.Sprintf("Failed to parse %q as an integer", parser.currentToken.Literal)
    parser.errors = append(parser.errors, msg)
    return nil
  }

  literal.Value = value
  return literal
}

func (parser *Parser) parseStringLiteral() ast.Expression {
  return &ast.StringLiteral{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIdentifier() ast.Expression {
  return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
  expression := &ast.PrefixExpression{
    Token: parser.currentToken,
    Operator: parser.currentToken.Literal,
  }

  parser.nextToken()

  expression.Right = parser.parseExpression(PREFIX)

  return expression
} // e.g. eval -5, PrefixExpression token is -, then advance the token,
  // then expression.Right is 5.


// Infix parser function
func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
  expression := &ast.InfixExpression{
    Token: parser.currentToken,
    Operator: parser.currentToken.Literal,
    Left: left,
  }

  precedence := parser.currentPrecendence()
  parser.nextToken()
  expression.Right = parser.parseExpression(precedence)

  return expression

}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
  expression := &ast.IndexExpression{
    Token: parser.currentToken, 
    Left: left,
  }
  parser.nextToken()
  expression.Index = parser.parseExpression(LOWEST)

  if !parser.expectPeek(token.RBRACKET) {
    return nil
  }

  return expression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
  expression := &ast.CallExpression{Token: parser.currentToken, Function: function}
  expression.Arguments = parser.parseExpressionList(token.RPAREN)
  return expression
}

func (parser *Parser) parseCallArguments() []ast.Expression {
  args := []ast.Expression{}

  if parser.peekTokenIs(token.RPAREN) {
    parser.nextToken()
    return args
  }

  parser.nextToken()
  args = append(args, parser.parseExpression(LOWEST))

  for parser.peekTokenIs(token.COMMA) {
    parser.nextToken()
    parser.nextToken()
    args = append(args, parser.parseExpression(LOWEST))
  }

  if !parser.expectPeek(token.RPAREN) {
    return nil
  }

  return args
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
  parser.nextToken()

  expression := parser.parseExpression(LOWEST)

  if !parser.expectPeek(token.RPAREN) {
    return nil
  }

  return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
  expression := &ast.IfExpression{Token: parser.currentToken}

  if !parser.expectPeek(token.LPAREN) {
    return nil
  }

  parser.nextToken()
  expression.Condition = parser.parseExpression(LOWEST)

  if !parser.expectPeek(token.RPAREN) {
    return nil
  }

  if !parser.expectPeek(token.LBRACE) {
    return nil
  }

  expression.Consequence = parser.parseBlockStatement()

  if parser.peekTokenIs(token.ELSE) {
    parser.nextToken()

    if !parser.expectPeek(token.LBRACE) {
      return nil
    }

    expression.Alternative = parser.parseBlockStatement()
  }

  return expression
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
  literal := &ast.FunctionLiteral{Token: parser.currentToken}

  if !parser.expectPeek(token.LPAREN) {
    return nil
  }
  
  literal.Parameters = parser.parseFunctionParameters()

  if !parser.expectPeek(token.LBRACE) {
    return nil
  }

  literal.Body = parser.parseBlockStatement()

  return literal
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
  block := &ast.BlockStatement{Token: parser.currentToken}
  block.Statements = []ast.Statement{}

  parser.nextToken()

  for !parser.currentTokenIs(token.RBRACE) && !parser.currentTokenIs(token.EOF) {
    statement := parser.parseStatement()
    if statement != nil {
      block.Statements = append(block.Statements, statement)
    }
    parser.nextToken()
  }

  return block
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
  identifiers := []*ast.Identifier{}

  if parser.peekTokenIs(token.RPAREN) {
    parser.nextToken()
    return identifiers
  }

  parser.nextToken()

  identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
  identifiers = append(identifiers, identifier)

  for parser.peekTokenIs(token.COMMA) {
    parser.nextToken()
    parser.nextToken()
    identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
    identifiers = append(identifiers, identifier)
  }

  if !parser.expectPeek(token.RPAREN) {
    return nil
  }

  return identifiers
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
  parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
  parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) peekPrecedence() int {
  if prec, ok := precedences[parser.peekToken.Type]; ok {
    return prec
  }

  return LOWEST
}

func (parser *Parser) currentPrecendence() int {
  if prec, ok := precedences[parser.currentToken.Type]; ok {
    return prec
  }

  return LOWEST
}

func (parser *Parser) noPrefixParseFnError(t token.TokenType) {
  msg := fmt.Sprintf("No prefix parser function found for %s", t)
  parser.errors = append(parser.errors, msg)
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

  parser.peekError(tokenType)
  return false
}

func (parser *Parser) peekError(tokenType token.TokenType) {
  msg := fmt.Sprintf("expected next token to be %s, but received %s", tokenType, parser.peekToken.Type)
  parser.errors = append(parser.errors, msg)
}

func (parser *Parser) nextToken() {
  parser.currentToken = parser.peekToken
  parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) parseBoolean() ast.Expression {
  return &ast.Boolean{Token: parser.currentToken, Value: parser.currentTokenIs(token.TRUE)}
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
  array := &ast.ArrayLiteral{Token: parser.currentToken}

  array.Elements = parser.parseExpressionList(token.RBRACKET)

  return array
}

func (parser *Parser) parseExpressionList(end token.TokenType) []ast.Expression { 
  list := []ast.Expression{}
  if parser.peekTokenIs(end) { 
    parser.nextToken()
    return list 
  }  
  
  parser.nextToken()
  list = append(list, parser.parseExpression(LOWEST))
  
  for parser.peekTokenIs(token.COMMA) { 
    parser.nextToken()
    parser.nextToken()
    list = append(list, parser.parseExpression(LOWEST))
  }

  if !parser.expectPeek(end) {
     return nil
  }

  return list
}

func (parser *Parser) parseHashLiteral() ast.Expression {
  hash := &ast.HashLiteral{Token: parser.currentToken}
  hash.Pairs = make(map[ast.Expression]ast.Expression)

  for !parser.peekTokenIs(token.RBRACE) {
    parser.nextToken()
    key := parser.parseExpression(LOWEST)

    if !parser.expectPeek(token.COLON) {
      return nil
    }

    parser.nextToken()
    value := parser.parseExpression(LOWEST)

    hash.Pairs[key] = value

    if !parser.peekTokenIs(token.RBRACE) && !parser.expectPeek(token.COMMA) {
      return nil
    }
  }

  if !parser.expectPeek(token.RBRACE) {
    return nil
  }

  return hash
}