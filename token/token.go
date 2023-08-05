package token

type TokenType string

type Token struct {
  Type    TokenType
  Literal string
}

const (
  ILLEGAL   = "ILLEGAL" // tokens/characters we don't know about/are undefined.
  EOF       = "EOF" // End of file

  // Variable identifiers and literal values
  IDENT     = "IDENT" // variable identifiers, foo, bar, user, etc.
  INT       = "INT"   // 123, 3555
  
  // Operators
  ASSIGN    = "="
  PLUS      = "+"
  MINUS     = "-"

  // Delimitters
  COMMA     = ","
  SEMICOLON = ";"
  LPAREN    = "("
  RPAREN    = ")"
  LBRACE    = "{"
  RBRACE    = "}"

  FUNCTION  = "FUNCTION"
  LET       = "LET"
)

var keywords = map[string]TokenType {
  "fn":     FUNCTION,
  "let":    LET,
}

func LookupIdent(ident string) TokenType {
  if tok, ok := keywords[ident]; ok {
    return tok 
  }
  return IDENT
}
