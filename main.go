package main

import (
	"log"

	"github.com/antonmedv/expr"
)

func main() {
	log.SetFlags(log.Lshortfile)
	p, err := expr.Compile("a > 1 ? a + 1 : 0", expr.AsInt(), expr.AllowUndefinedVariables())
	if err != nil {
		log.Println(err)
		return
	}

	r, err := expr.Run(p, map[string]interface{}{
		"a": 1000,
	})
	log.Printf("r=%+v err=%+v", r, err)
}
