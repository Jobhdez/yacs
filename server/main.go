package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ScmLike struct {
	Exp string `json:"exp"`
}

type Result struct {
	Exp string `json:"exp"`
}

func Compile(exp string) string {
	ast, _ := Parse(exp)

	monAst := ToAnf(ast)
	ss := ToSelect(monAst)
	return ToAssembly(ss)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusNotFound)
	body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", body)
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func CompileHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	// Handle OPTIONS method for CORS preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var exp ScmLike
	err = json.Unmarshal(d, &exp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	selectIns := Compile(exp.Exp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Result{Exp: selectIns})
}

func main() {
	var PORT = ":1234"
	if len(os.Args) > 1 {
		PORT = ":" + os.Args[1]
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/api/compiler", http.HandlerFunc(CompileHandler))
	mux.HandleFunc("/", defaultHandler)
	log.Println("listening on port:", PORT)
	log.Fatal(s.ListenAndServe())
}
