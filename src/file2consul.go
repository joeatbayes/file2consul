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
	"path/filepath"
)

// Parse a simple key value properties file
// adding the values and keys found to the target
// dictionary.  Using a supplied dictionary because
// the intent is to call it with multiple files and
// interpolate the results.
func loadFileAsDict(inFiName string, target map[string]string, pargs *jutil.ParsedCommandArgs, doInterpolate bool) {
	start := jutil.Nowms()
	verboseFlg := pargs.Exists("verbose")
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("ERROR: loadFileAsDict opening input file ", inFiName, " err=", err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	printLinesFlg := pargs.Exists("printlines")
	if verboseFlg {
		fmt.Println("loadFileAsDict ", inFiName, "doInterpolate=", doInterpolate)
	}
	lineCnt := 0
	lastKey := ""
	for scanner.Scan() {
		aline := scanner.Text()
		lineCnt++
		aline = strings.TrimSpace(aline)
		if printLinesFlg {
			fmt.Println("#", lineCnt, "\t", aline)
		}
		if len(aline) <= 0 {
			continue
		}

		if strings.HasPrefix(aline, "#") {
			continue
		}

		if strings.HasPrefix(aline, "+") && (len(aline) > 1) && (lastKey > "") {
			appendVal := strings.TrimSpace(aline[1:])
			target[lastKey] = target[lastKey] + "\t" + pargs.Interpolate(appendVal)
			continue
		}

		arr := strings.SplitN(aline, "=", 2)
		if len(arr) != 2 {
			fmt.Println("NOTE: line#", lineCnt, "fails split on = test", " line=", aline)
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
		target[aKey] = aVal
		if strings.HasPrefix(aVal, "@") {
			//  If first character of the values starts with @
			//  then attempt to load the contents of that file
			//  and replace the value with those contents.
			currDir := filepath.Dir(inFiName)
			reqPath := aVal[1:]
			newPath := currDir + "/" + reqPath
			if verboseFlg {
				fmt.Println("Attempt to load referenced File requested=", aVal, " derivedName=", newPath)
			}
			if jutil.IsDirectory(newPath) {
				fmt.Println("ERROR: loadFileAsDict: Included File can not be directory requested=", newPath, " err=", err)
			} else if jutil.Exists(newPath) {
				dat, err := ioutil.ReadFile(newPath)
				if err != nil {
					fmt.Println("ERROR: loadFileAsDict: Error reading referenced file=", newPath, " err=", err)
				} else {
					aVal = string(dat)
					if verboseFlg {
						fmt.Println("loadFileAsDict: Loaded ", len(dat), " bytes from referenced file=", newPath)
					}
					aVal = pargs.Interpolate(aVal) // want the fully interpolated value override for the print.
					target[aKey] = aVal
				}
			} else {
				fmt.Println("ERROR: loadFileAsDict Can not find referenced file requested=", aVal, " derivedName=", newPath)
			}
		} // include file

		lastKey = aKey
		if verboseFlg {
			fmt.Println("key=", aKey, " val=", aVal)
		}
	}
	if verboseFlg {
		jutil.Elap("loadFileAsDict "+inFiName, start, jutil.Nowms())
	}
}

func loadInFiles(inPaths []string, pargs *jutil.ParsedCommandArgs, doInterpolate bool) map[string]string {
	verboseFlg := pargs.Exists("verbose")
	start := jutil.Nowms()
	inDict := make(map[string]string) // Stores out built up set of config parameters
	// Load all the specified paths into the dictionary
	for i, apath := range inPaths {
		if verboseFlg {
			fmt.Println("Process Path=", apath, "ndx=", i)
		}
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
				if verboseFlg {
					fmt.Println("  Process File=", fName)
				}
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

// Create log file entry for the current run showing variables saved to consul
// along with details of the run.
func logChangedItems(sdict map[string]string, fiName string, pargs *jutil.ParsedCommandArgs) {
	//verboseFlg := pargs.Exists("verbose")
}

func saveValuesToConsul(sdict map[string]string, serverURI string, pargs *jutil.ParsedCommandArgs) {
	verboseFlg := pargs.Exists("verbose")
	dc := pargs.Sval("dc", "")
	start := jutil.Nowms()
	if verboseFlg {
		fmt.Println("saveValuesToConsul ", serverURI, " numItems=", len(sdict))
	}
	for akey, aval := range sdict {
		jutil.SetConsulKey(serverURI, akey, aval, dc)
	}

	jutil.Elap("Save Values to Consul "+serverURI, start, jutil.Nowms())

}

func saveValuesToConsuls(sdict map[string]string, serverURIs []string, pargs *jutil.ParsedCommandArgs) {
	for _, serverURI := range serverURIs {
		saveValuesToConsul(sdict, serverURI, pargs)
	}
}

func main() {
	start := jutil.Nowms()
	args := os.Args
	pargs := jutil.ParseCommandLine(args)

	if len(args) < 2 { // || (pargs.Exists("h")) || (pargs.Exists("help")) {
		msg := `EXAMPLE
  file2consul -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -IN=data/config/simple/template::data/config/simple/prod::data/config/simple/uat::data/config/simple/joes-dev.prop.txt -uri=http://127.0.0.1:8500 -CACHE=data/{env}.CACHE.b64.txt
 
  
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
	
   -PATHDELIM =  Delimiter to use when splitting list of 
     files, paths or URI for fields like -URI, -IN.   Defaults
	 to :: when not set.
	       
   -PRINTLINES when this value is specified the
     system will print every input line as it is read
     to help in diagnostics.
     
   -VERBOSE When this value is specified the system 
     will print additional details about values as 
     they are set or re-set during the run. 
	
   -DC  when set this value is passed to the data 
     center attribute in the PUT call to consul. 
		
   -appname = variable used for interpolation
   -env =  variable used for interpolation
   -company = variable used for interpolation
   -appname = varabile used for interpolation

   other named parameters are treated in interpolated values
   Most common error is forgetting - as prefix for named parms			
		 -`
		fmt.Println(msg)
	}

	pathDelim := pargs.Sval("pathdelim", "::")
	inPaths := strings.Split(pargs.Sval("in", "data/config/simple/basic"), pathDelim)
	serverURIStr := pargs.Sval("uri", "http://127.0.0.1:8500") // need the basic string to support the NONE check
	serverURIs := strings.Split(serverURIStr, pathDelim)
	serverURIStrFlgChk := strings.ToUpper(serverURIStr)
	cacheFiName := pargs.Interpolate(pargs.Sval("cache", ""))
	verboseFlg := pargs.Exists("verbose")
	saveReadable := pargs.Interpolate(pargs.Sval("savereadable", "NONE"))
	if verboseFlg {
		fmt.Println("pathDelim=", pathDelim, " inPaths=", inPaths, " consul server URI=", serverURIs, " cacheFiName=", cacheFiName)
	}
	inDict := loadInFiles(inPaths, pargs, true)
	if verboseFlg {
		fmt.Println("inDict=", inDict)
	}

	if serverURIStrFlgChk != "NONE" {
		fmt.Println("NOTE: Skip save to Consule and update of Cache file because -uri == NONE")
	}

	if cacheFiName > "" {
		// If the cache file is specified then we only want
		// to send Delta to consul.
		if jutil.Exists(cacheFiName) {
			cacheDict := jutil.LoadDictFile(cacheFiName, true, pargs)
			deltaDict := jutil.CompareStrDict(cacheDict, inDict)
			if len(deltaDict) < 1 {
				fmt.Println("NOTE: No Values have changed since cache last updated, Not need to update consul")
			} else if serverURIStrFlgChk != "NONE" {
				if verboseFlg {
					fmt.Println("NOTE: Only need to save ", len(deltaDict), " items to Consul due to cache check")
				}
				saveValuesToConsuls(deltaDict, serverURIs, pargs)
			}
		} else if serverURIStrFlgChk != "NONE" {
			// dictionary file did not exist yet so need to
			// send entire set to consul.
			saveValuesToConsuls(inDict, serverURIs, pargs)
		}
		if serverURIStrFlgChk != "NONE" {
			jutil.SaveDictToFile(inDict, cacheFiName, true, pargs)
		}
	} else if serverURIStrFlgChk != "NONE" {
		// NO Cache is specified so no need to save the
		// values to consul
		// save entire file set to Consul since
		// we are not using the cache.
		saveValuesToConsuls(inDict, serverURIs, pargs)
	}

	if saveReadable != "NONE" {
		// save a human readable file to allow easier debugging
		jutil.SaveDictToFile(inDict, saveReadable, false, pargs)
	}
	jutil.Elap("file2Consul complete run", start, jutil.Nowms())

	keys := jutil.GetConsulKeys(serverURIs[0], "", "+", "")
	fmt.Println(" keys=", keys)
} //main
