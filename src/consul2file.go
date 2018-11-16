/* consoleSaveKeys.go
Save the list of keys in the current Consul
// NOTE:  Not sure how many keys can be returned in a single call.

// TODO: TEST with very large key set. If a problem then Modify
//  to walk by level so we never get more than one folder at a time.
*/
package main

import (
	//"bufio"
	"fmt"

	//"io/ioutil"
	"jutil"

	//"net/http"
	"os"
	//"strings"
	//"time"
	//"net/url"
	//"path/filepath"
)

func main() {
	start := jutil.Nowms()
	args := os.Args
	pargs := jutil.ParseCommandLine(args)

	if len(args) < 2 { // || (pargs.Exists("h")) || (pargs.Exists("help")) {
		msg := `EXAMPLE
  consul2file -OUT=data/consul-save-key.txt -uri=http://127.0.0.1:8500 -prefix=/
    Saves consul server contents to a local file.
   
   -OUT=name of file to write contents into. 

   -URI=uri to reach consul server.   
        If seprated by ; will save to each server listed
		defaults to http://127.0.0.1:8500 if not specified
		
   -PREFIX = prefix of the key to start the extract. When 
     not specified uses no prefix or starts at root.
	
   -DELIM =  Delimiter to use when splitting path to provide
     folder style functionality.  Defaults to / when 
	 not set. 
	
   KEYSONLY = when specified system will fetch values
      from consul and include them in the saved file.

   URLEncode = When specified any values will be URL escaped 
	
   B64 = When specified any values written will be base64 encoded.
    unless URLEncode flag is set. 
	
   -TODO: MAKETREE = When specified will use the out name as a 
     directory and will create a directory tree reflecting the 
	 containment folders. Uses / as the default delimiter.

	       
		 -`
		fmt.Println(msg)
	}

	delim := pargs.Sval("delim", "/")
	prefix := pargs.Sval("prefix", "")
	out := pargs.Sval("out", "data/save-keys.txt")
	serverURI := pargs.Sval("uri", "http://127.0.0.1:8500") // need the basic string to support the NONE check
	verboseFlg := pargs.Exists("verbose")
	b64Flg := pargs.Exists("b64")
	//URLEncode := pargs.Exists("urlencode")

	if verboseFlg {
		fmt.Println("out=", out, " consul server URI=", serverURI, " prefix=", prefix, " delim=", delim)
	}
	if prefix == "/" {
		prefix = "" // map single slash to empty to make URI formatting easier
	}

	keys := jutil.GetConsulKeys(serverURI, "", "+", "")
	fmt.Println(" keys=", keys)
	if pargs.Exists("keysonly") {
		jutil.SaveStrsToFile(keys, out, false, pargs)
	} else {

		tout := make(map[string]string)
		for ndx, key := range keys {
			val, err := jutil.GetConsulVal(serverURI, key, "")
			if err != nil {
				fmt.Printf("ndx=%d err=%s key=%s val=%s\n", err, ndx, key, val)
			}

			fmt.Printf("ndx=%d key=%s val=%s\n", ndx, key, val)
			// Parse value input with JSON
			tout[key] = val
			//if b64Val {
			//saveStr = key + "=" + base64.StdEncoding.EncodeToString([]byte(val)) + "\n"
			//}
		}

		jutil.SaveDictToFile(tout, out, b64Flg, pargs)
	}
	// TODO:  Fetch the value for each Key from consul
	//  and write the actual values to the save file.
	jutil.Elap("consul2file complete run", start, jutil.Nowms())
} //main
