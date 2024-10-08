package main

import (
	"fmt"
	"strings"
)

type MonExpression interface{}

type MonLet struct {
	MonBindings []MonBinding
	Body        MonExpression
}

type MonBinding struct {
	Name  string
	Value MonExpression
}

type MonVar struct {
	Name string
}

type MonInt struct {
	Value int
}

func ToAnf(expr Expression) MonExpression {
	switch e := expr.(type) {
	case IntLiteral:
		return MonInt{Value: e.Value}
	case LetExpr:
		monBindings := make([]MonBinding, len(e.Bindings))
		for i, bind := range e.Bindings {
			monBindings[i] = MonBinding{
				Name:  bind.Name,
				Value: ToAnf(bind.Value),
			}
		}
		monBody := ToAnf(e.Body)
		return MonLet{MonBindings: monBindings, Body: monBody}
	case Var:
		return MonVar{Name: e.Name}
	default:
		return nil
	}
}

func PrintMonLet(letExpr MonLet) {
	fmt.Println("MonLet Expression:")
	for _, binding := range letExpr.MonBindings {
		fmt.Printf("Binding: %s = ", binding.Name)
		PrintMon(binding.Value)
	}
	fmt.Printf("Body: ")
	PrintMon(letExpr.Body)
	fmt.Println()
}

func PrintMon(mon MonExpression) {
	switch e := mon.(type) {
	case MonInt:
		fmt.Printf("MonInt(%d)\n", e.Value)
	case MonVar:
		fmt.Printf("MonVar(%s)\n", e.Name)
	case MonLet:
		PrintMonLet(e)
	default:
		fmt.Println("Unknown MonExpression")
	}
}

func Parse(input string) (Expression, error) {
	lexer := &Lexer{}
	lexer.Init(strings.NewReader(input))

	if yyParse(lexer) == 0 {
		return lexer.result, nil
	} else {
		return nil, fmt.Errorf("Parsing error")
	}
}

func main() {
	input := "(let ((x 3)) 4)"
	ast, err := Parse(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	monAst := ToAnf(ast)
	PrintMon(monAst)
}
