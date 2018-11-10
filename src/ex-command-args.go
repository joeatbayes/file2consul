// ex-command-args.go
// Example of parsing different kinds of command args
//  by Joe Ellsworth 2017-07-21
//  license: MIT
//  I provide consulting services http://BayesAnalytic.com/contact

package main

import (
	"fmt"
	"jutil"
	"os"
)

func main() {
	args := os.Args

	fmt.Println("ex-command-args")
	fmt.Println("tTest with complex command patterns ")
	fmt.Printf("\tShow Args by Index\n")
	numArgs := len(os.Args)
	for ndx := 0; ndx < numArgs; ndx++ {
		fmt.Printf("\t\tndx=%v  arg=%v\n", ndx, args[ndx])
	}
	fmt.Println("\n")

	// SEE:   jutil.CommandLineParserTest()
	//    for sample of easy ways to interact with
	//    parsed command line aurguments
	jutil.CommandLineParserTest()

	// Most simple use of Command Line Parser
	fmt.Println("\nDemonsrate Simple Use of Command Parser")
	pargs := jutil.ParseCommandLine(args)

	// Dump as easy read string
	fmt.Printf("parsed Args\n%s\n", pargs)

	// Get string parm with default
	classFileName := pargs.Sval("class", "../../sampleDefaultFiName.txt")
	fmt.Printf("classFileName=%s", classFileName)

}
