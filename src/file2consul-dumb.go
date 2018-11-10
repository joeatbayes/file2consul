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

		msg := `

  file2consul-dumb -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -FILE=data/config/simple/template/basic.prop.txt -uri=http://127.0.0.1:8500
 
  
   -file=name of input paramter file
   -uri=uri to reach console server
   -appname = variable used for interpolation
   -env =  variable used for interpolation
   -company = variable used for interpolation
   -appname = varabile used for interpolation
 
   other variables can be defined as needed
   variables are not case sensitive.
		-`
		fmt.Println(msg)
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
		aKey := pargs.Interpolate(arr[0])
		aVal := pargs.Interpolate(arr[1])
		fmt.Println("after interpolate aKey=", aKey, " aVal=", aVal)
		jutil.SetConsulKey(serverURI, aKey, aVal)
		fmt.Println("\n")

	} // for
} //main
