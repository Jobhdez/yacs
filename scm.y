%{
package main
import (
	"fmt"
	"text/scanner"
	"strconv"
)

// Define the structure for expressions.
type Expression interface{}

type IntLiteral struct {
	Value int
}

type Var struct {
	Name string
}

type Token struct {
    token int
    literal string
}

type IfExpr struct {
	Cond Expression
	Then Expression
	Else Expression
}

type LetExpr struct {
	Bindings []Binding
	Body     Expression
}

type Binding struct {
	Name  string
	Value Expression
}

type Application struct {
	Func Expression
	Args []Expression
}

type DefineExpr struct {
	Name  string
	Value Expression
}

// Function to print a LetExpr
func PrintLetExpr(letExpr LetExpr) {
	fmt.Println("Let Expression:")
	for _, binding := range letExpr.Bindings {
		fmt.Printf("Binding: %s = ", binding.Name)
		PrintExpr(binding.Value)
	}
	fmt.Printf("Body: ")
	PrintExpr(letExpr.Body)
	fmt.Println()
}

// Function to print any Expression
func PrintExpr(expr Expression) {
	switch e := expr.(type) {
	case IntLiteral:
		fmt.Printf("IntLiteral(%d)\n", e.Value)
	case Var:
		fmt.Printf("Var(%s)\n", e.Name)
	case IfExpr:
		fmt.Printf("IfExpr(Cond: ")
		PrintExpr(e.Cond)
		fmt.Printf(" Then: ")
		PrintExpr(e.Then)
		fmt.Printf(" Else: ")
		PrintExpr(e.Else)
		fmt.Println(")")
	case LetExpr:
		PrintLetExpr(e)
	case Application:
		fmt.Printf("Application(Func: ")
		PrintExpr(e.Func)
		fmt.Printf(" Args: ")
		for _, arg := range e.Args {
			PrintExpr(arg)
		}
		fmt.Println(")")
	case DefineExpr:
		fmt.Printf("DefineExpr(Name: %s, Value: ", e.Name)
		PrintExpr(e.Value)
		fmt.Println(")")
	default:
		fmt.Println("Unknown expression")
	}
}
%}

%union {
	token  Token
	expr   Expression
	str    string
	intval int
}

%type<expr> program
%type<expr> expr
%type<expr> binding
%type<expr> expr_list

%token<str> NAME
%token<intval> INTEGER
%token LPAREN RPAREN

%token<str> LET
%token IF DEFINE LAMBDA EQ

%start program

%%

program:
	expr {
		$$ = $1
		 yylex.(*Lexer).result = $$
		  //PrintExpr($$)
	}

expr:
	INTEGER {
		$$ = IntLiteral{Value: $1}
	}
	| NAME {
		$$ = Var{Name: $1}
	}
	| LPAREN LET LPAREN binding RPAREN expr RPAREN {
		$$ = LetExpr{Bindings: []Binding{$4.(Binding)}, Body: $6}
	}
	| LPAREN IF expr expr expr RPAREN {
		$$ = IfExpr{Cond: $3, Then: $4, Else: $5}
	}
	| LPAREN DEFINE NAME expr RPAREN {
		$$ = DefineExpr{Name: $3, Value: $4}
	}
	| LPAREN expr expr_list RPAREN {
		$$ = Application{Func: $2, Args: $3.([]Expression)}
	}

binding:
	LPAREN NAME expr RPAREN {
		$$ = Binding{Name: $2, Value: $3}
	}

expr_list:
	/* empty */ {
		$$ = []Expression{}
	}
	| expr expr_list {
		$$ = append([]Expression{$1}, $2.([]Expression)...)
	}

%%

type Lexer struct {
	scanner.Scanner
	result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok := l.Scan()
	lit := l.TokenText()

	switch tok {
	case scanner.Int:
		// Handle integers
		tokVal, _ := strconv.Atoi(lit)
		lval.intval = tokVal
		return INTEGER
	case '(':
		// Handle left parenthesis
		return LPAREN
	case ')':
		// Handle right parenthesis
		return RPAREN
	case scanner.Ident:
		// Handle identifiers like 'if', 'let', 'define', or variables
		switch lit {
		case "if":
			return IF
		case "let":
			return LET
		case "define":
			return DEFINE
		default:
			// Any other identifier is treated as a variable name
			lval.str = lit
			return NAME
		}
	case '=':
		return EQ
	}

	// If no valid token was found, return 0 (end of input or error)
	return 0
}

func (l *Lexer) Error(e string) {
	fmt.Printf("Lex error: %s\n", e)
}
