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
  STRING    = "STRING"
  
  // Operators
  ASSIGN    = "="
  PLUS      = "+"
  MINUS     = "-"
  BANG      = "!"
  ASTERISK  = "*"
  SLASH     = "/"
  LT        = "<"
  GT        = ">"

  // Delimitters
  COMMA     = ","
  SEMICOLON = ";"
  LPAREN    = "("
  RPAREN    = ")"
  LBRACE    = "{"
  RBRACE    = "}"
  EQ        = "=="
  NOT_EQ    = "!="

  // Keywords
  FUNCTION  = "FUNCTION"
  LET       = "LET"
  TRUE      = "TRUE"
  FALSE     = "FALSE"
  IF        = "IF"
  ELSE      = "ELSE"
  RETURN    = "RETURN"
)

var keywords = map[string]TokenType {
  "fn":     FUNCTION,
  "let":    LET,
  "true":   TRUE,
  "false":  FALSE,
  "if":     IF,
  "else":   ELSE,
  "return": RETURN,
}

func LookupIdent(ident string) TokenType {
  if tok, ok := keywords[ident]; ok {
    return tok 
  }
  return IDENT
}
