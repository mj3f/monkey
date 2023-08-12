package ast

import (
  "monkey/token"
)

type Node interface {
  TokenLiteral() string
}

type Statement interface {
  Node
  statementNode()
}

type Expression interface {
  Node
  expressionNode()
}

type Program struct { // The root node of the AST the parser produces.
  Statements []Statement // slice of statement nodes, (Every valid program written in Monkey is a series of statements)
}

func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  }

  return ""
}

type LetStatement struct {
  Name *Identifier
  Value Expression
  Token token.Token // token.LET
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
  return ls.Token.Literal
}

type ReturnStatement struct {
  Token token.Token
  ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
  return rs.Token.Literal
}

type Identifier struct {
  Value string
  Token token.Token // token.IDENT
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
  return i.Token.Literal
}
