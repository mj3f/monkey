package lexer

import (
	"monkey/token"
)

type Lexer struct {
  input           string
  position        int // currnet position in input (current character)
  readPosition    int // current reading position in input (after current char)
  ch              byte // current character under evaluation
}

func New(input string) *Lexer {
  lexer := &Lexer{input: input}
  lexer.readChar()
  return lexer
}

// Reads the next character in the input sequence, sets the lexer to 0 if we reach the end.
func (lexer *Lexer) readChar() {
  if lexer.readPosition >= len(lexer.input) {
    lexer.ch = 0
  } else {
    lexer.ch = lexer.input[lexer.readPosition]
  }

  lexer.position = lexer.readPosition
  lexer.readPosition += 1
}

func (lexer *Lexer) NextToken() token.Token {
  var tok token.Token
  lexer.skipWhitespace()

  switch lexer.ch  {
    case '=':
      if lexer.peekChar() == '=' {
        ch := lexer.ch
        lexer.readChar()
        literal := string(ch) + string(lexer.ch)
        tok = token.Token{Type: token.EQ, Literal: literal}
      } else {
        tok = newToken(token.ASSIGN,lexer.ch) 
      }
    case ';':
      tok = newToken(token.SEMICOLON, lexer.ch) 
    case '(':
      tok = newToken(token.LPAREN, lexer.ch) 
    case ')':
      tok = newToken(token.RPAREN, lexer.ch) 
    case ',':
      tok = newToken(token.COMMA, lexer.ch) 
    case '+':
      tok = newToken(token.PLUS, lexer.ch) 
    case '{':
      tok = newToken(token.LBRACE, lexer.ch) 
    case '}':
      tok = newToken(token.RBRACE, lexer.ch)
    case '!':
      if lexer.peekChar() == '=' {
        ch := lexer.ch
        lexer.readChar()
        literal := string(ch) + string(lexer.ch)
        tok = token.Token{Type: token.NOT_EQ, Literal: literal}
      } else {
        tok = newToken(token.BANG, lexer.ch)
      }
    case '/':
      tok = newToken(token.SLASH, lexer.ch)
    case '*':
      tok = newToken(token.ASTERISK, lexer.ch)
    case '<':
      tok = newToken(token.LT, lexer.ch)
    case '>':
      tok = newToken(token.GT, lexer.ch)
    case 0:
      tok.Literal = ""
      tok.Type = token.EOF
    default:
      if isLetter(lexer.ch) {
        tok.Literal = lexer.readIdentifier()
        tok.Type = token.LookupIdent(tok.Literal)
        return tok // readIdentifier calls readChar repeatedly to update positions, so we early return here.
      } else if isDigit(lexer.ch) {
        tok.Type = token.INT
        tok.Literal = lexer.readNumber()
        return tok
      } else {
        tok = newToken(token.ILLEGAL, lexer.ch)
      }
  }

  lexer.readChar()
  return tok
}

func (lexer *Lexer) readIdentifier() string {
  position := lexer.position
  for isLetter(lexer.ch) {
    lexer.readChar()
  }
  return lexer.input[position:lexer.position] // create slice containing just the identifier of the variable.
}

func isLetter(ch byte) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?' // allow names like foo_bar and name?
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
  return token.Token{Type: tokenType, Literal: string(ch)}
}

func (lexer *Lexer) skipWhitespace() {
  for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
    lexer.readChar()
  }
}

func (lexer *Lexer) readNumber() string {
  position := lexer.position
  for isDigit(lexer.ch) {
    lexer.readChar()
  }
  return lexer.input[position:lexer.position]
}

func isDigit(ch byte) bool {
  return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) peekChar() byte {
  if lexer.readPosition >= len(lexer.input) {
    return 0
  }

  return lexer.input[lexer.readPosition]
}
