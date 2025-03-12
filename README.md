# This is an interpreter for 'my own' programming language, called 'Monkey'!.

I purchased a copy of [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball, this repo is re-creating the code shown in the book step-by-step to write my own version of the interpreter.

It is fully unit tested, to run any of the tests, i.e. `go test ./parser`.

To run this, you'll need to install Go, afterwards just run the following command to launch the REPL: `go run main.go`.

## Examples

Once the REPL is launched, try some of these examples:

```monkey
let x = 5;
let people = { "name": "John", "age": 20 };

puts(people["name"]);

let y = x * 2;

puts(y);
```
