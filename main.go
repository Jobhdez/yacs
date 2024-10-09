package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Mon Intermediate

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

type MonIf struct {
	Cond MonExpression
	Then MonExpression
	Else MonExpression
}

type MonBinary struct {
	Op    string
	Left  MonExpression
	Right MonExpression
}

type MonWhile struct {
	Cnd  MonExpression
	Body MonExpression
}

type MonSet struct {
	Var string
	Exp MonExpression
}

// Select Instructor
type Instructions struct {
	Instructs [][]string
}

func PrintLetExpr(letExpr MonLet) {
	fmt.Println("Let Expression:")
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
		PrintLetExpr(e)
	case MonIf:
		fmt.Printf("MonIfExpr(Cond: ")
		PrintMon(e.Cond)
		fmt.Printf(" Then: ")
		PrintMon(e.Then)
		fmt.Printf(" Else: ")
		PrintMon(e.Else)
		fmt.Println(")")
	case MonBinary:
		fmt.Printf("MonBinaryOp(Operator: %s, Left: ", e.Op)
		PrintMon(e.Left)
		fmt.Printf(", Right: ")
		PrintMon(e.Right)
		fmt.Println(")")
	case MonWhile:
		fmt.Printf("MonWhile:")
		PrintMon(e.Cnd)
		PrintMon(e.Body)
	case MonSet:
		PrintMon(e.Var)
		PrintMon(e.Exp)
	default:
		fmt.Println("Unknown MonExpression")
	}
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
	case IfExpr:
		cnd := ToAnf(e.Cond)
		thn := ToAnf(e.Then)
		els := ToAnf(e.Else)
		return MonIf{Cond: cnd, Then: thn, Else: els}
	case BinaryOp:
		left := ToAnf(e.Left)
		right := ToAnf(e.Right)
		return MonBinary{Op: e.Operator, Left: left, Right: right}
	case WhileExpr:
		cnd := ToAnf(e.Cnd)
		body := ToAnf(e.Body)
		return MonWhile{Cnd: cnd, Body: body}
	case SetExpr:
		variable := e.Name
		exp := ToAnf(e.Value)
		return MonSet{Var: variable, Exp: exp}
	default:
		return nil
	}
}

func SelectInstructions(expr MonExpression) Instructions {
	switch e := expr.(type) {
	case MonInt:
		instructions := make([][]string, 0)
		strnum := strconv.Itoa(e.Value)
		movinstruction := []string{"movq", strnum, "%rdi"}
		callinstruction := []string{"callq", "print_int"}

		instructions = append(instructions, movinstruction, callinstruction)
		return Instructions{Instructs: instructions}

	case MonVar:
		instructions := make([][]string, 0)
		movinstruction := []string{"movq", e.Name, "%rdi"}
		callinstruction := []string{"callq", "print_int"}

		instructions = append(instructions, movinstruction, callinstruction)
		return Instructions{Instructs: instructions}

	case MonLet:
		instructions := make([][]string, 0)
		binding := e.MonBindings[0]

		switch val := binding.Value.(type) {
		case MonInt:
			strnum := strconv.Itoa(val.Value)
			movinstruction := []string{"movq", strnum, binding.Name}
			instructions = append(instructions, movinstruction)
			bodyinstructions := SelectInstructions(e.Body)
			instructions = append(instructions, bodyinstructions.Instructs...)
			return Instructions{Instructs: instructions}
		default:
			fmt.Println("Unsupported MonExpression in Let")
			return Instructions{Instructs: [][]string{}}
		}

	case MonIf:
		condInstructions := SelectInstructions(e.Cond)
		thenInstructions := SelectInstructions(e.Then)
		elseInstructions := SelectInstructions(e.Else)
		instructions := append(condInstructions.Instructs, thenInstructions.Instructs...)
		instructions = append(instructions, elseInstructions.Instructs...)
		return Instructions{Instructs: instructions}

	case MonBinary:
		op := e.Op
		switch op {
		case "<":
			instructions := make([][]string, 0)

			rightExpr, ok := e.Right.(MonInt)
			if !ok {
				fmt.Println("Expected MonInt for the right side of binary operation")
				return Instructions{Instructs: [][]string{}}
			}
			leftExpr, ok := e.Left.(MonVar)
			if !ok {
				fmt.Println("Expected MonVar for the left side of binary operation")
				return Instructions{Instructs: [][]string{}}
			}

			strnum := strconv.Itoa(rightExpr.Value)
			cmpin := []string{"cmpq", strnum, leftExpr.Name}
			instructions = append(instructions, cmpin)

			return Instructions{Instructs: instructions}

		default:
			fmt.Println("Unsupported binary operator")
			return Instructions{Instructs: [][]string{}}
		}
	case MonWhile:
		cnd := SelectInstructions(e.Cnd)
		cndins := cnd.Instructs
		body := SelectInstructions(e.Body)
		bodyins := body.Instructs

		instructions := make([][]string, 0)
		jllabel := [][]string{{"label", "loop"}}
		jlbody := append(jllabel, bodyins...)
		ins := append(instructions, jlbody...)
		jlin := [][]string{{"jl", "loop"}}
		inscmp := append(ins, cndins...)
		cmpjl := append(inscmp, jlin...)

		return Instructions{Instructs: cmpjl}
	case MonSet:
		switch exp := e.Exp.(type) {
		case MonInt:
			strnum := strconv.Itoa(exp.Value)
			ins := [][]string{{"movq", strnum, e.Var}}
			return Instructions{Instructs: ins}
		case MonBinary:
			er := exp.Right
			var strnum string
			switch exp2 := er.(type) {
			case MonInt:
				strnum = strconv.Itoa(exp2.Value)
				ins := [][]string{{"movq", strnum, "%rax"}, {"addq", "%rax", e.Var}}
				return Instructions{Instructs: ins}
			default:
				fmt.Println("Unsupported right-hand expression for binary operator")
			}
		default:
			fmt.Println("Unsupported expression in MonSet")
		}
		return Instructions{Instructs: [][]string{}}

	default:
		return Instructions{Instructs: [][]string{}}
	}
}

func PrintSelect(ins Instructions) {
	fmt.Println(ins.Instructs)
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
	//input := "(let ((i 0)) (if (< i 3) 3 5))"
	//input := "(let ((i 0)) (while (< i 5) i))"
	input := "(set i (+ i 1))"
	
	ast, err := Parse(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	PrintMon(ast)
	monAst := ToAnf(ast)
	PrintMon(monAst)
	ss := SelectInstructions(monAst)
	PrintSelect(ss)
}
