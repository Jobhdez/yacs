clean:
	rm scheme_parser y.output
yacc:
	 goyacc -o scheme_parser.go scheme_parser.y

main:
	go build -o scheme_parser main.go compiler.go scheme_parser.go




