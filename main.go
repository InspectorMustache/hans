package main

import (
	"fmt"
)

func main() {
	ccs := new(charInfo)
	ccs.populateStruct("        一       1       一       一       1               *       0               M       *")
	fmt.Println(ccs)
}
