package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "
const MONKEY_FACE = `            
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)
  env := object.NewEnvironment()

  for {
    fmt.Fprintf(out, PROMPT)
    scanned := scanner.Scan()
    if !scanned {
      return 
    }

    line := scanner.Text()
    lexer := lexer.New(line)
    parser := parser.New(lexer)

    program := parser.ParseProgram()

    if len(parser.Errors()) != 0 {
      printParserErrors(out, parser.Errors())
      continue
    }

    evaluated := evaluator.Eval(program, env)

    if evaluated != nil {
      io.WriteString(out, evaluated.Inspect())
      io.WriteString(out, "\n")
    }
  }
}

func printParserErrors(out io.Writer, errors []string) {
  io.WriteString(out, MONKEY_FACE)
  io.WriteString(out, "Woops! We ran into some monkey business here!\n")
  io.WriteString(out, " errors during parsing:\n")
  for _, msg := range errors {
    io.WriteString(out, "\t" + msg + "\n")
  }
}
