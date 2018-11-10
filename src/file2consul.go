/* file2console-dumb.go
Most simple / dumb version of saving lines from specified input file to consul.
Not recomended but it can help validate consul is responding correctly
*/
package main

import (
	"bufio"
	"fmt"

	"io/ioutil"
	"jutil"

	//"net/http"
	"os"
	"strings"
	//"time"
)

// Parse a simple key value properties file
// adding the values and keys found to the target
// dictionary.  Using a supplied dictionary because
// the intent is to call it with multiple files and
// interpolate the results.
func loadFileAsDict(inFiName string, target map[string]string, pargs *jutil.ParsedCommandArgs, doInterpolate bool) {
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
		var aKey string
		var aVal string
		if doInterpolate {
			aKey = pargs.Interpolate(arr[0])
			aVal = pargs.Interpolate(arr[1])
		} else {
			aKey = arr[0]
			aVal = arr[1]
		}
		// TODO:  Detect @ as first character of file to load file relative to current file
		//   for value instead of using included string.
		target[aKey] = aVal
		fmt.Println("key=", aKey, " val=", aVal)
	}
}

func loadInFiles(inPaths []string, pargs *jutil.ParsedCommandArgs, doInterpolate bool) map[string]string {
	inDict := make(map[string]string) // Stores out built up set of config parameters
	// Load all the specified paths into the dictionary
	for i, apath := range inPaths {
		fmt.Println("Process Path=", apath, "ndx=", i)
		if !(jutil.Exists(apath)) {
			fmt.Println("ERROR: input path does not exist", apath)
			continue
		}

		if jutil.IsDirectory(apath) {
			files, err := ioutil.ReadDir(apath)
			if err != nil {
				fmt.Println("ERROR: could not read dir ", apath, " error=", err)
				continue
			}

			for _, f := range files {
				fName := apath + "/" + f.Name()
				fmt.Println("  Process File=", fName)
				loadFileAsDict(fName, inDict, pargs, doInterpolate)
			}
		} else {
			// process path as simple input file
			loadFileAsDict(apath, inDict, pargs, doInterpolate)
		}
	} // for paths
	return inDict
}

func saveValuesToConsul(sdict map[string]string, serverURI string) {
	for akey, aval := range sdict {
		jutil.SetConsulKey(serverURI, akey, aval)
	}
}

func saveValuesToConsuls(sdict map[string]string, serverURIs []string) {
	for _, serverURI := range serverURIs {
		saveValuesToConsul(sdict, serverURI)
	}
}

func main() {

	args := os.Args
	pargs := jutil.ParseCommandLine(args)

	if len(args) < 2 { // || (pargs.Exists("h")) || (pargs.Exists("help")) {
		msg := `EXAMPLE
  file2consul -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -IN=data/config/simple/template;data/config/simple/prod;data/config/simple/uat;data/config/simple/joes-dev.prop.txt -uri=http://127.0.0.1:8500
 
  
   -IN=name of input paramter file or directory
       If named resource is directory will process all 
	   files in that directory.    Multiple inputs
	   can be specified seprated by ;.  Each input set
	   will be processed in order with any duplicate 
	   keys in subsequent files overriding those 
	   previously defined for the same key.   This 
	   provides a simple form of inheritance where
	   only the values that change between enviornment
	   need to be listed while the rest can be inherited
	   from a common parent.  If not specified defaults
	   to data/config/simple/basic
	  

   -URI=uri to reach console server.   
        If seprated by ; will save to each server listed
		defaults to http://127.0.0.1:8500 if not specified
		
   -appname = variable used for interpolation
   -env =  variable used for interpolation
   -company = variable used for interpolation
   -appname = varabile used for interpolation

   other named parameters are treated in interpolated values
   Most common error is forgetting - as prefix for named parms			
		 -`
		fmt.Println(msg)
	}

	inPaths := strings.Split(pargs.Sval("in", "data/config/simple/basic"), ";")
	serverURIs := strings.Split(pargs.Sval("uri", "http://127.0.0.1:8500"), ";")
	fmt.Println("inPaths=", inPaths, " consul server URI=", serverURIs)

	inDict := loadInFiles(inPaths, pargs, true)
	fmt.Println("inDict=", inDict)
	saveValuesToConsuls(inDict, serverURIs)

} //main
