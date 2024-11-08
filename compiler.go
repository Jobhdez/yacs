package main

import (
	"fmt"
	"strconv"
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

type MonBegin struct {
	Exps []MonExpression
}

// Select Instructions
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
		fmt.Printf("MonSet(Var: %s, Exp: ", e.Var)
		PrintMon(e.Exp)
		fmt.Println(")")
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
	case BeginExpr:
		exps := make([]MonExpression, len(e.Exprs))
		for i := range e.Exprs {
			exps[i] = ToAnf(e.Exprs[i])
		}
		return MonBegin{Exps: exps}
	default:
		return nil
	}
}

func ToSelect(expr MonExpression) Instructions {
	stack := make(map[string]string)
	var counter = 0
	instructions, stack_ := SelectInstructions(expr, stack, counter)
	if ((stack_ % 16) == 0) {
		stack_ = stack_ + 0
	} else {
		stack_ = stack_ + 8
	}
	var stack_n string
	stack_n = strconv.Itoa(stack_)
	prelude := [][]string{
		{"\t.globl ", "main\n"},
		{"main:\n"},
		{"\tpushq ", "%rbp\n"},
		{"\tmovq ", "%rsp, ", "%rbp\n"},
		{"\tsubq ", "$" + stack_n + ", ", "%rsp\n"}}
	
	conclusion := [][]string{
		{"\nconclusion:\n"},
		{"\taddq ", "$"+ stack_n + ", ", "%rsp\n"},
		{"\tpopq ", "%rbp\n"},
		{"\tretq"}}
	
	instructions_ := append(prelude, instructions.Instructs...)
	instructions_ = append(instructions_, conclusion...)

	return Instructions{Instructs: instructions_}
}

func SelectInstructions(expr MonExpression, stack map[string]string, counter int) (Instructions, int) {
	switch e := expr.(type) {
	case MonInt:
		instructions := [][]string{
			{"\tmovq ", "$" + strconv.Itoa(e.Value) + ", ", "%rdi\n"},
			{"\tcallq ", "print_int\n"},
		}
		return Instructions{Instructs: instructions}, counter

	case MonVar:
		movInstruction := []string{"\tmovq ", stack[e.Name] + ", ", "%rdi\n"}
		callInstruction := []string{"\tcallq ", "print_int\n"}
		instructions := [][]string{movInstruction, callInstruction}
		return Instructions{Instructs: instructions}, counter

	case MonLet:
		instructions := make([][]string, 0)
		binding := e.MonBindings[0]

		switch val := binding.Value.(type) {
		case MonInt:
			strnum := strconv.Itoa(val.Value)
			if _, ok := stack[binding.Name]; ok {
				movInstruction := []string{"\tmovq ", strnum + ", ", stack[binding.Name] + "\n"}
				instructions = append(instructions, movInstruction)
			} else {
				counter += 8
				stackLocation := "-" + strconv.Itoa(counter) + "(%rbp)"
				stack[binding.Name] = stackLocation
				movInstruction := []string{"\tmovq ", "$" + strnum +", ", stackLocation + "\n"}
				instructions = append(instructions, movInstruction)
			}
			bodyInstructions, _ := SelectInstructions(e.Body, stack, counter)
			instructions = append(instructions, bodyInstructions.Instructs...)
			return Instructions{Instructs: instructions}, counter
		default:
			fmt.Println("Unsupported MonExpression in Let")
			return Instructions{Instructs: [][]string{}}, counter
		}
	case MonBegin:
		instructions := make([][]string, 0)
		for _, exp := range e.Exps {
			expInstructions, _ := SelectInstructions(exp, stack, counter)
			instructions = append(instructions, expInstructions.Instructs...)
		}
		return Instructions{Instructs: instructions}, counter

	case MonIf:
		condInstructions, _ := SelectInstructions(e.Cond, stack, counter)
		thenInstructions, _ := SelectInstructions(e.Then, stack, counter)
		elseInstructions, _ := SelectInstructions(e.Else, stack, counter)

		instructions := append(condInstructions.Instructs, []string{"\tjl ", "label1\n"})
		instructions = append(instructions, []string{"\tjmp ", "label2\n"})
		instructions = append(instructions, []string{"\nlabel1:\n"})
		instructions = append(instructions, thenInstructions.Instructs...)
		instructions = append(instructions, []string{"\tjmp ", "conclusion\n"})
		instructions = append(instructions, []string{"\nlabel2:\n"})
		instructions = append(instructions, elseInstructions.Instructs...)
		instructions = append(instructions, []string{"\tjmp ", "conclusion\n"})
		return Instructions{Instructs: instructions}, counter
	case MonSet:
		switch ins := e.Exp.(type) {
		case MonBinary:
			instructions, _ := SelectInstructions(ins, stack, counter)
			return instructions, counter
		default:
			fmt.Println("unsupported set expr")
			return Instructions{Instructs: [][]string{}}, counter
		}
	case MonWhile:
		cndInstructions, _ := SelectInstructions(e.Cnd, stack, counter)
		bodyInstructions, _ := SelectInstructions(e.Body, stack, counter)
		instructions := append([][]string{{"\nloop_99:\n"}}, bodyInstructions.Instructs...)
		instructions = append(instructions, cndInstructions.Instructs...)
		instructions = append(instructions, []string{"\tjl ", "loop_99\n"})
		return Instructions{Instructs: instructions}, counter
		
	case MonBinary:
	switch e.Op {
	case "<":
		rightVal, ok := e.Right.(MonInt)
		leftName, leftExists := stack[e.Left.(MonVar).Name]
		if !ok || !leftExists {
			fmt.Println("Unsupported or missing MonExpression in Binary Op")
			return Instructions{Instructs: [][]string{}}, counter
		}
		n := strconv.Itoa(rightVal.Value)
		instructions := [][]string{
			{"\tcmpq ", "$" + n + ", ", leftName + "\n"},
		}
		return Instructions{Instructs: instructions}, counter
	case "+":
		rightVal, ok := e.Right.(MonInt)
		leftName, leftExists := stack[e.Left.(MonVar).Name]
		if !ok || !leftExists {
			fmt.Println("Unsupported or missing MonExpression in Binary Op")
			return Instructions{Instructs: [][]string{}}, counter
		}
		n := strconv.Itoa(rightVal.Value)
		instructions := [][]string{
			{"\tmovq ", "$" + n + ", ",  "%rax\n"},
			{"\taddq ", "%rax ", leftName + "\n"},
		}
		return Instructions{Instructs: instructions}, counter
	default:
		fmt.Println("Unsupported operator in MonBinary")
		return Instructions{Instructs: [][]string{}}, counter
	}

	default:
		fmt.Println("Unsupported expression type")
		return Instructions{Instructs: [][]string{}}, counter
	}
}

func PrintSelect(ins Instructions) {
	for _, instr := range ins.Instructs {
		fmt.Println(instr)
	}
}

func InstructionsToString(ins Instructions) string {
	var result string
	for _, instr := range ins.Instructs {
		if (len(instr) == 1) {
			result += fmt.Sprintf("%s", instr[0])
		} else if (len(instr) == 2) {
			result += fmt.Sprintf("%s %s", instr[0], instr[1])
		} else {
			result += fmt.Sprintf("%s %s %s", instr[0], instr[1], instr[2])
		}
	}
	return result
}

/*
func main() {
	// input := "(if (< 2 3) 2 3)"
	//input := "(let ((sum 0)) (let ((i 0)) (if (< sum 4) i 9)))"
	input := "(let ((sum 0)) (let ((i 0)) (begin (while (< i 5) (begin (set sum (+ sum 3)) (set i (+ i 1)))) sum)))"
	ast, _ := Parse(input) // Parse function needs to be defined

	monAst := ToAnf(ast)
	PrintMon(monAst)
	ss := ToSelect(monAst)
	fmt.Println(InstructionsToString(ss))
}
*/
