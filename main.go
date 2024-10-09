package main

import (
	"fmt"
)

func main() {
	//input := "(let ((i 0)) (if (< i 3) 3 5))"
	//input := "(let ((i 0)) (while (< i 5) i))"
	//input := "(set i (+ i 1))"
	input := "(let ((i 0)) (begin (while (< i 4) (set i (+ i 1))) i))"
	
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
