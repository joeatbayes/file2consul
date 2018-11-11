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
		aKey = strings.TrimSpace(aKey)
		aVal = strings.TrimSpace(aVal)
		// TODO:  Detect @ as first character of file to load file relative to current file
		//   for value instead of using included string.
		target[aKey] = aVal
		fmt.Println("key=", aKey, " val=", aVal)
	}
}

func loadInFiles(inPaths []string, pargs *jutil.ParsedCommandArgs, doInterpolate bool) map[string]string {
	start := jutil.Nowms()
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
	jutil.Elap("Load all files", start, jutil.Nowms())
	return inDict
}

func obtainChangedItems(sdict map[string]string) map[string]string {

	return nil
}

// Create log file entry for the current run showing variables saved to consul
// along with details of the run.
func logChangedItems(sdict map[string]string, fiName string, pargs *jutil.ParsedCommandArgs) {

}

func saveValuesToConsul(sdict map[string]string, serverURI string) {
	start := jutil.Nowms()
	for akey, aval := range sdict {
		jutil.SetConsulKey(serverURI, akey, aval)
	}
	jutil.Elap("Save Values to Consul "+serverURI, start, jutil.Nowms())
}

func saveValuesToConsuls(sdict map[string]string, serverURIs []string) {
	for _, serverURI := range serverURIs {
		saveValuesToConsul(sdict, serverURI)
	}
}

func main() {
	start := jutil.Nowms()
	args := os.Args
	pargs := jutil.ParseCommandLine(args)

	if len(args) < 2 { // || (pargs.Exists("h")) || (pargs.Exists("help")) {
		msg := `EXAMPLE
  file2consul -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -IN=data/config/simple/template;data/config/simple/prod;data/config/simple/uat;data/config/simple/joes-dev.prop.txt -uri=http://127.0.0.1:8500 -CACHE=data/{env}.CACHE.b64.txt
 
  
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
	  

   -URI=uri to reach consul server.   
        If seprated by ; will save to each server listed
		defaults to http://127.0.0.1:8500 if not specified
		
   -CACHE = name of files to use as cache file.  This file is 
     read and compared to the post processing set to determine
	 what values need to be saved to consul.  It is also re-written
	 and end of run when defined.  If you want to clear cache
	 delete the file before running the utility.  This value
	 is subjected to interpolation so you can use things like
	 enviornment as part of file name.
		
   -appname = variable used for interpolation
   -env =  variable used for interpolation
   -company = variable used for interpolation
   -appname = varabile used for interpolation

   other named parameters are treated in interpolated values
   Most common error is forgetting - as prefix for named parms			
		 -`
		fmt.Println(msg)
	}

	inPaths := strings.Split(pargs.Sval("in", "data/config/simple/basic"), "::")
	serverURIs := strings.Split(pargs.Sval("uri", "http://127.0.0.1:8500"), "::")
	cacheFiName := pargs.Interpolate(pargs.Sval("cache", ""))
	fmt.Println("inPaths=", inPaths, " consul server URI=", serverURIs, "cacheFiName="+cacheFiName)
	inDict := loadInFiles(inPaths, pargs, true)
	//fmt.Println("inDict=", inDict)

	if cacheFiName > "" {
		// If the cache file is specified then we only want
		// to send Delta to consul.
		if jutil.Exists(cacheFiName) {
			cacheDict := jutil.LoadDictFile(cacheFiName, true)
			deltaDict := jutil.CompareStrDict(cacheDict, inDict)
			if len(deltaDict) < 1 {
				fmt.Println("NOTE: No Values have changed, Not need to update consul")
			} else {
				saveValuesToConsuls(deltaDict, serverURIs)
			}
		} else {
			// dictionary file did not exist yet so need to
			// send entire set to consul.
			saveValuesToConsuls(inDict, serverURIs)
		}
		jutil.SaveDictToFile(inDict, cacheFiName, true)
	} else {
		// save entire file set to Consul since
		// we are not using the cache.
		saveValuesToConsuls(inDict, serverURIs)
	}

	if cacheFiName > "" {

	}
	jutil.Elap("file2Consul complete run", start, jutil.Nowms())

} //main
