/* file2console-dumb.go
Most simple / dumb version of saving lines from specified input file to consul.
Not recomended but it can help validate consul is responding correctly
*/
package main

import (
	"bufio"
	"fmt"
	//"io/ioutil"
	"jutil"
	//"net/http"
	"os"
	"strings"
	//"time"
)

func main() {

	args := os.Args
	pargs := jutil.ParseCommandLine(args)

	if len(args) < 2 { // || (pargs.Exists("h")) || (pargs.Exists("help")) {
		fmt.Println("file2consul-dumb\n")
		fmt.Println("\texample:  executableName -uri=http://127.0.0.1:8500  -inFile=data/simple-config/basic.prxxxop.txt -ENV=LOCAL")
		fmt.Println("\t  -uri= uri of the consul server to save the values in.")
		fmt.Println("\t     If seprated by , can save to more than one server")
		fmt.Println("\t     defaults to http://127.0.0.1:8500 if not specified")
		fmt.Println("\t   -inFile= Name of input file to process if not specified")
		fmt.Println("\t      then defaults to data/simple-config/basic.prop.txt")
		fmt.Println("\t   -ENV=LOCAL enviornment variable that can be interpolated into output")
		fmt.Println("\t      defaults to LOCAL if not set.")
		fmt.Println("\t  other named parameters are treated in interpolated values")
		fmt.Println("\t  Most common error is forgetting - as prefix for named parms")
		fmt.Println("\t   Number of Args: %v  args: %s\n", len(args), args)
	}

	inFiName := pargs.Sval("file", "data/simple-config/basic.prop.txt")
	serverURI := pargs.Sval("uri", "http://127.0.0.1:8500")
	fmt.Println("inFile=", inFiName, " consul server URI=", serverURI)
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("error opening ID file ", inFiName, " err=", err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	lineCnt := 0
	for scanner.Scan() {
		aline := scanner.Text()
		lineCnt++
		if len(aline) <= 0 {
			continue
		}
		//fmt.Println(aline)
		aline = strings.TrimSpace(aline)
		if strings.HasPrefix(aline, "#") {
			continue
		}

		arr := strings.SplitN(aline, "=", 2)
		if len(arr) != 2 {
			fmt.Println("line#", lineCnt, "fails split on = test", " line=", aline)
			continue
		}
		aKey := jutil.Interpolate(arr[0], pargs)
		aVal := jutil.Interpolate(arr[1], pargs)
		fmt.Println("after interpolate aKey=", aKey, " aVal=", aVal)
		jutil.SetConsulKey(serverURI, aKey, aVal)
		fmt.Println("\n")

	} // for
} //main
