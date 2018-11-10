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
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var str = "Old Value!"
var ParmMatch, ParmErr = regexp.Compile("\\{.*?\\}")

func setConsulKey(serverURI string, key string, val string) {
	fmt.Println("setConsulKey key=", key, " value=", val)
	uri := serverURI + "/v1/kv/" + key
	val = "sample data value"
	//fmt.Println("uri=", uri)
	start := time.Now().UTC()
	hc := http.Client{}
	req, err := http.NewRequest("PUT", uri, strings.NewReader(val))
	if err != nil {
		fmt.Println("Error: opening uri=", uri, " err=", err, " key=", key, "  val=", val)
		return
	}

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("Error: uri=", uri, " err=", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error: Expected Status Code 200, Got ", resp.StatusCode)
	} else {
		//fmt.Println("Sucess:")
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	jutil.TimeTrack(os.Stdout, start, "finished single PUT uri="+uri+"\n")
	//fmt.Println("statusCode=", resp.StatusCode)
	//fmt.Println(" body=", string(body))
}

func interpolate(str string, parg *jutil.ParsedCommandArgs) string {
	ms := ParmMatch.FindAllIndex([]byte(str), -1)
	if len(ms) < 1 {
		return str // no match found
	}
	//sb := strings.Builder
	var sb []string
	last := 0
	slen := len(str)
	for _, m := range ms {
		start, end := m[0]+1, m[1]-1
		//fmt.Printf("m[0]=%d m[1]=%d match = %q\n", m[0], m[1], str[start:end])
		if start > last-1 {
			// add the string before the match to the buffer
			sb = append(sb, str[last:start-1])
		}
		aMatchStr := strings.ToLower(str[start:end])
		// substitute match string with parms value
		// or add it back in with the {} protecting it
		// TODO: Add lookup from enviornment variable
		//  if do not find it in the command line parms
		lookVal := parg.Sval(strings.ToLower(aMatchStr), "{"+aMatchStr+"}")
		//fmt.Printf("matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
		sb = append(sb, lookVal)
		last = end + 1
	}
	if last < slen-1 {
		// append any remaining characters after
		// end of the last match
		sb = append(sb, str[last:slen])
	}
	return strings.Join(sb, "")
}

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
		aKey := interpolate(arr[0], pargs)
		aVal := interpolate(arr[1], pargs)
		fmt.Println("after interpolate aKey=", aKey, " aVal=", aVal)
		setConsulKey(serverURI, aKey, aVal)
		fmt.Println("\n")

	} // for
} //main
