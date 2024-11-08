clean:
	rm service y.output compiler 
yacc:
	 goyacc -o parser.go parser.y

main:
	go build -o service main.go compiler.go parser.go


compile:
	go build -o compiler compiler.go parser.go

